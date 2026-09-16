package vectorxws

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const onboardingMaxBody = 2 << 20

var onboardingBackend = "http://127.0.0.1:8070"

var onboardingActions = map[string]string{
	"init":         "/api-ble/init",
	"scan":         "/api-ble/scan",
	"connect":      "/api-ble/connect",
	"pin":          "/api-ble/send_pin",
	"wifi-status":  "/api-ble/get_wifi_status",
	"wifi-scan":    "/api-ble/scan_wifi",
	"wifi-connect": "/api-ble/connect_wifi",
	"authenticate": "/api-ble/do_auth",
	"disconnect":   "/api-ble/disconnect",
	"oskr-setup":   "/api-ssh/setup",
	"oskr-status":  "/api-ssh/get_setup_status",
}

var onboardingHTTPClient = &http.Client{Timeout: 90 * time.Second}

func onboardingAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	if !sameOrigin(r) {
		http.Error(w, "request origin not allowed", http.StatusForbidden)
		return
	}
	action := strings.TrimPrefix(r.URL.Path, "/api/onboarding/")
	upstreamPath, ok := onboardingActions[action]
	if !ok {
		http.Error(w, "unknown onboarding action", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost && !(r.Method == http.MethodGet && action == "oskr-status") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	base, err := url.Parse(onboardingBackend)
	if err != nil {
		http.Error(w, "onboarding service is not configured", http.StatusServiceUnavailable)
		return
	}
	base.Path = upstreamPath
	base.RawQuery = r.URL.RawQuery
	body := http.MaxBytesReader(w, r.Body, onboardingMaxBody)
	request, err := http.NewRequestWithContext(r.Context(), r.Method, base.String(), body)
	if err != nil {
		http.Error(w, "could not prepare onboarding request", http.StatusBadRequest)
		return
	}
	request.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	response, err := onboardingHTTPClient.Do(request)
	if err != nil {
		http.Error(w, "VectorX onboarding service is unavailable", http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	w.WriteHeader(response.StatusCode)
	_, _ = io.Copy(w, io.LimitReader(response.Body, onboardingMaxBody))
}

func sameOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	parsed, err := url.Parse(origin)
	return err == nil && strings.EqualFold(parsed.Host, r.Host)
}

func onboardingError(raw string) error {
	value := strings.TrimSpace(raw)
	if value == "" || strings.HasPrefix(strings.ToLower(value), "error") {
		return fmt.Errorf("%s", value)
	}
	return nil
}
