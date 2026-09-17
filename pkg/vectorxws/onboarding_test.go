package vectorxws

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSameOrigin(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{name: "browser same origin", origin: "http://vectorx.local:8070", want: true},
		{name: "non browser client", origin: "", want: true},
		{name: "foreign site", origin: "https://example.com", want: false},
		{name: "lookalike host", origin: "http://vectorx.local:8070.example.com", want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8070/api/onboarding/scan", nil)
			r.Host = "vectorx.local:8070"
			if test.origin != "" {
				r.Header.Set("Origin", test.origin)
			}
			if got := sameOrigin(r); got != test.want {
				t.Fatalf("sameOrigin() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRemoveRobotPreservesOtherRobotsAndRestartsService(t *testing.T) {
	dir := t.TempDir()
	registry := filepath.Join(dir, "botSdkInfo.json")
	original := `{"global_guid":"global","robots":[{"esn":"005070ac","ip_address":"192.168.1.34","guid":"secret-a","activated":true},{"esn":"00112233","ip_address":"192.168.1.35","guid":"secret-b","activated":true}]}`
	if err := os.WriteFile(registry, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}
	oldControl := wirePodServiceControl
	var actions []string
	wirePodServiceControl = func(action string) error { actions = append(actions, action); return nil }
	defer func() { wirePodServiceControl = oldControl }()

	r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8070/api/remove_robot", strings.NewReader("esn=005070ac"))
	r.Host = "vectorx.local:8070"
	r.Header.Set("Origin", "http://vectorx.local:8070")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	removeRobotHandler(w, r, registry)
	if w.Code != http.StatusOK {
		t.Fatalf("response = %d %q", w.Code, w.Body.String())
	}
	updated, err := os.ReadFile(registry)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(updated), "005070ac") || !strings.Contains(string(updated), "00112233") {
		t.Fatalf("unexpected registry: %s", updated)
	}
	if strings.Join(actions, ",") != "stop,start" {
		t.Fatalf("service actions = %v", actions)
	}
	backups, _ := filepath.Glob(registry + ".bak-vectorx-remove-*")
	if len(backups) != 1 {
		t.Fatalf("backup count = %d", len(backups))
	}
}

func TestUpdateRobotIPPreservesOtherRobotsAndRestartsService(t *testing.T) {
	dir := t.TempDir()
	registry := filepath.Join(dir, "botSdkInfo.json")
	original := `{"global_guid":"global","robots":[{"esn":"005070ac","ip_address":"192.168.1.34","guid":"secret-a","activated":true},{"esn":"00112233","ip_address":"192.168.1.35","guid":"secret-b","activated":true}]}`
	if err := os.WriteFile(registry, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}
	oldControl := wirePodServiceControl
	var actions []string
	wirePodServiceControl = func(action string) error { actions = append(actions, action); return nil }
	defer func() { wirePodServiceControl = oldControl }()

	r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8070/api/update_robot_ip", strings.NewReader("esn=005070ac&ip=192.168.1.17"))
	r.Host = "vectorx.local:8070"
	r.Header.Set("Origin", "http://vectorx.local:8070")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	updateRobotIPHandler(w, r, registry)
	if w.Code != http.StatusOK {
		t.Fatalf("response = %d %q", w.Code, w.Body.String())
	}
	updated, err := os.ReadFile(registry)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(updated), `"ip_address": "192.168.1.17"`) || !strings.Contains(string(updated), `"ip_address": "192.168.1.35"`) {
		t.Fatalf("unexpected registry: %s", updated)
	}
	if strings.Join(actions, ",") != "stop,start" {
		t.Fatalf("service actions = %v", actions)
	}
	backups, _ := filepath.Glob(registry + ".bak-vectorx-ip-*")
	if len(backups) != 1 {
		t.Fatalf("backup count = %d", len(backups))
	}
}

func TestUpdateRobotIPRejectsInvalidAddress(t *testing.T) {
	dir := t.TempDir()
	registry := filepath.Join(dir, "botSdkInfo.json")
	if err := os.WriteFile(registry, []byte(`{"global_guid":"global","robots":[{"esn":"005070ac","ip_address":"192.168.1.34","guid":"secret-a","activated":true}]}`), 0644); err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8070/api/update_robot_ip", strings.NewReader("esn=005070ac&ip=127.0.0.1"))
	r.Host = "vectorx.local:8070"
	r.Header.Set("Origin", "http://vectorx.local:8070")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	updateRobotIPHandler(w, r, registry)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("response = %d %q", w.Code, w.Body.String())
	}
}

func TestOnboardingProxyUsesAllowlistAndPreservesFormBody(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api-ble/send_pin" {
			t.Fatalf("upstream path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "pin=123456" {
			t.Fatalf("upstream body = %q", body)
		}
		_, _ = w.Write([]byte("success"))
	}))
	defer upstream.Close()

	oldBackend, oldClient := onboardingBackend, onboardingHTTPClient
	onboardingBackend, onboardingHTTPClient = upstream.URL, upstream.Client()
	defer func() { onboardingBackend, onboardingHTTPClient = oldBackend, oldClient }()

	r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8070/api/onboarding/pin", strings.NewReader("pin=123456"))
	r.Host = "vectorx.local:8070"
	r.Header.Set("Origin", "http://vectorx.local:8070")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	onboardingAPIHandler(w, r)
	if w.Code != http.StatusOK || strings.TrimSpace(w.Body.String()) != "success" {
		t.Fatalf("response = %d %q", w.Code, w.Body.String())
	}
}

func TestOnboardingProxyRejectsUnknownActionAndForeignOrigin(t *testing.T) {
	for _, test := range []struct {
		path   string
		origin string
		code   int
	}{
		{path: "/api/onboarding/shell", origin: "http://vectorx.local:8070", code: http.StatusNotFound},
		{path: "/api/onboarding/scan", origin: "https://example.com", code: http.StatusForbidden},
	} {
		r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8070"+test.path, nil)
		r.Host = "vectorx.local:8070"
		r.Header.Set("Origin", test.origin)
		w := httptest.NewRecorder()
		onboardingAPIHandler(w, r)
		if w.Code != test.code {
			t.Fatalf("%s returned %d, want %d", test.path, w.Code, test.code)
		}
	}
}

func TestRobotSettingsProxyUsesAllowlist(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api-sdk/volume" {
			t.Fatalf("upstream path = %q", r.URL.Path)
		}
		_ = r.ParseForm()
		if r.Form.Get("serial") != "005070ac" || r.Form.Get("volume") != "4" {
			t.Fatalf("upstream form = %v", r.Form)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer upstream.Close()
	oldBackend, oldClient := onboardingBackend, onboardingHTTPClient
	onboardingBackend, onboardingHTTPClient = upstream.URL, upstream.Client()
	defer func() { onboardingBackend, onboardingHTTPClient = oldBackend, oldClient }()
	r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8070/api/robot-settings/volume", strings.NewReader("serial=005070ac&volume=4"))
	r.Host = "vectorx.local:8070"
	r.Header.Set("Origin", "http://vectorx.local:8070")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	robotSettingsAPIHandler(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("response = %d %q", w.Code, w.Body.String())
	}
}

func TestRobotSettingsRejectsUnknownAction(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8070/api/robot-settings/shell", strings.NewReader("serial=005070ac"))
	w := httptest.NewRecorder()
	robotSettingsAPIHandler(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("response = %d", w.Code)
	}
}

func TestRobotSettingsRejectsUnapprovedQuickAction(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8080/api/robot-settings/quick-action", strings.NewReader("serial=005070ac&intent=run_shell"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	robotSettingsAPIHandler(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("response = %d %q", w.Code, w.Body.String())
	}
}

func TestRobotSettingsPreservesPhotoContentType(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api-sdk/get_image_thumb" {
			t.Fatalf("upstream path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte{0xff, 0xd8, 0xff, 0xd9})
	}))
	defer upstream.Close()
	oldBackend, oldClient := onboardingBackend, onboardingHTTPClient
	onboardingBackend, onboardingHTTPClient = upstream.URL, upstream.Client()
	defer func() { onboardingBackend, onboardingHTTPClient = oldBackend, oldClient }()
	r := httptest.NewRequest(http.MethodPost, "http://vectorx.local:8080/api/robot-settings/photo-thumb", strings.NewReader("serial=005070ac&id=1"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	robotSettingsAPIHandler(w, r)
	if w.Code != http.StatusOK || w.Header().Get("Content-Type") != "image/jpeg" {
		t.Fatalf("response = %d content-type=%q", w.Code, w.Header().Get("Content-Type"))
	}
}

func TestCheckAvailableUpdateComparesLocalAndRemoteRevisions(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	oldHome := os.Getenv("VECTORX_HOME")
	oldGit := gitCommandOutput
	os.Setenv("VECTORX_HOME", dir)
	defer func() {
		os.Setenv("VECTORX_HOME", oldHome)
		gitCommandOutput = oldGit
	}()

	gitCommandOutput = func(_ string, args ...string) ([]byte, error) {
		key := strings.Join(args, " ")
		switch key {
		case "fetch --no-tags --prune origin":
			return []byte(""), nil
		case "rev-parse --verify @{u}":
			return []byte("remote\n"), nil
		case "show @{u}:pkg/vectorxws/webserver.go":
			return []byte("package vectorxws\n\nconst VECTORX_VERSION = \"RELEASE_26\"\n"), nil
		}
		return nil, errors.New("unexpected git command: " + key)
	}

	status := checkAvailableUpdate()
	if status.Result != "ok" || !status.UpdateAvailable || status.AvailableVersion != "RELEASE_26" {
		t.Fatalf("unexpected status: %+v", status)
	}
}
