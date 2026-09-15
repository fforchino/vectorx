package vectorxws

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

var robotSettingsActions = map[string]string{
	"get":           "/api-sdk/get_sdk_settings",
	"volume":        "/api-sdk/volume",
	"locale":        "/api-sdk/locale",
	"eye-color":     "/api-sdk/eye_color",
	"location":      "/api-sdk/location",
	"timezone":      "/api-sdk/timezone",
	"time-12":       "/api-sdk/time_format_12",
	"time-24":       "/api-sdk/time_format_24",
	"temp-c":        "/api-sdk/temp_c",
	"temp-f":        "/api-sdk/temp_f",
	"button-vector": "/api-sdk/button_hey_vector",
	"button-alexa":  "/api-sdk/button_alexa",
}

func robotSettingsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	if !sameOrigin(r) {
		http.Error(w, `{"error":"request origin not allowed"}`, http.StatusForbidden)
		return
	}
	action := strings.TrimPrefix(r.URL.Path, "/api/robot-settings/")
	upstreamPath, ok := robotSettingsActions[action]
	if !ok {
		http.Error(w, `{"error":"unknown robot setting"}`, http.StatusNotFound)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, onboardingMaxBody)
	if err := r.ParseForm(); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	serial := strings.ToLower(strings.TrimSpace(r.Form.Get("serial")))
	if len(serial) != 8 {
		http.Error(w, `{"error":"invalid robot identifier"}`, http.StatusBadRequest)
		return
	}
	for _, c := range serial {
		if !strings.ContainsRune("0123456789abcdef", c) {
			http.Error(w, `{"error":"invalid robot identifier"}`, http.StatusBadRequest)
			return
		}
	}

	base, err := url.Parse(onboardingBackend)
	if err != nil {
		http.Error(w, `{"error":"robot service is not configured"}`, http.StatusServiceUnavailable)
		return
	}
	base.Path = upstreamPath
	values := url.Values{}
	for key, entries := range r.Form {
		for _, value := range entries {
			values.Add(key, value)
		}
	}
	request, err := http.NewRequestWithContext(r.Context(), http.MethodPost, base.String(), strings.NewReader(values.Encode()))
	if err != nil {
		http.Error(w, `{"error":"could not prepare robot request"}`, http.StatusBadRequest)
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := onboardingHTTPClient.Do(request)
	if err != nil {
		http.Error(w, `{"error":"robot service is unavailable"}`, http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	w.WriteHeader(response.StatusCode)
	_, _ = io.Copy(w, io.LimitReader(response.Body, onboardingMaxBody))
}
