package vectorxws

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const robotSettingsMaxBody = 12 << 20

type robotSettingsAction struct {
	path   string
	fields map[string]func(string) bool
}

func oneOf(values ...string) func(string) bool {
	allowed := make(map[string]bool, len(values))
	for _, value := range values {
		allowed[value] = true
	}
	return func(value string) bool { return allowed[value] }
}

func integerBetween(min, max int) func(string) bool {
	return func(value string) bool { n, err := strconv.Atoi(value); return err == nil && n >= min && n <= max }
}

func decimalBetween(min, max float64) func(string) bool {
	return func(value string) bool {
		n, err := strconv.ParseFloat(value, 64)
		return err == nil && n >= min && n <= max
	}
}

func nonEmptyLimited(max int) func(string) bool {
	return func(value string) bool {
		value = strings.TrimSpace(value)
		return value != "" && len(value) <= max && !strings.ContainsAny(value, "\r\n\x00")
	}
}

var robotSettingsActions = map[string]robotSettingsAction{
	"get":              {path: "/api-sdk/get_sdk_settings"},
	"volume":           {path: "/api-sdk/volume", fields: map[string]func(string) bool{"volume": integerBetween(0, 5)}},
	"locale":           {path: "/api-sdk/locale", fields: map[string]func(string) bool{"locale": oneOf("en-US", "en-GB", "en-AU", "de-DE", "fr-FR", "ja-JP")}},
	"eye-color":        {path: "/api-sdk/eye_color", fields: map[string]func(string) bool{"color": integerBetween(0, 6)}},
	"custom-eye-color": {path: "/api-sdk/custom_eye_color", fields: map[string]func(string) bool{"hue": decimalBetween(0, 1), "sat": decimalBetween(0, 1)}},
	"location":         {path: "/api-sdk/location", fields: map[string]func(string) bool{"location": nonEmptyLimited(120)}},
	"timezone":         {path: "/api-sdk/timezone", fields: map[string]func(string) bool{"timezone": nonEmptyLimited(100)}},
	"time-12":          {path: "/api-sdk/time_format_12"},
	"time-24":          {path: "/api-sdk/time_format_24"},
	"temp-c":           {path: "/api-sdk/temp_c"},
	"temp-f":           {path: "/api-sdk/temp_f"},
	"button-vector":    {path: "/api-sdk/button_hey_vector"},
	"button-alexa":     {path: "/api-sdk/button_alexa"},
	"behavior-assume":  {path: "/api-sdk/assume_behavior_control", fields: map[string]func(string) bool{"priority": oneOf("default", "high", "override")}},
	"behavior-release": {path: "/api-sdk/release_behavior_control"},
	"say-text":         {path: "/api-sdk/say_text", fields: map[string]func(string) bool{"text": nonEmptyLimited(599)}},
	"move-wheels":      {path: "/api-sdk/move_wheels", fields: map[string]func(string) bool{"lw": integerBetween(-250, 250), "rw": integerBetween(-250, 250)}},
	"move-lift":        {path: "/api-sdk/move_lift", fields: map[string]func(string) bool{"speed": integerBetween(-3, 3)}},
	"move-head":        {path: "/api-sdk/move_head", fields: map[string]func(string) bool{"speed": integerBetween(-3, 3)}},
	"mirror":           {path: "/api-sdk/mirror_mode", fields: map[string]func(string) bool{"enable": oneOf("true", "false")}},
	"camera-frame":     {path: "/api-sdk/camera_frame"},
	"camera-stop":      {path: "/api-sdk/stop_cam_stream"},
	"alexa-sign-in":    {path: "/api-sdk/alexa_sign_in"},
	"alexa-sign-out":   {path: "/api-sdk/alexa_sign_out"},
	"quick-action":     {path: "/api-sdk/cloud_intent", fields: map[string]func(string) bool{"intent": oneOf("explore_start", "intent_imperative_dance", "intent_system_sleep", "intent_imperative_fetchcube", "intent_system_charger")}},
	"faces":            {path: "/api-sdk/get_faces"},
	"face-rename":      {path: "/api-sdk/rename_face", fields: map[string]func(string) bool{"id": integerBetween(0, 2147483647), "oldname": nonEmptyLimited(100), "newname": nonEmptyLimited(100)}},
	"face-delete":      {path: "/api-sdk/delete_face", fields: map[string]func(string) bool{"id": integerBetween(0, 2147483647)}},
	"stim-start":       {path: "/api-sdk/begin_event_stream"},
	"stim-stop":        {path: "/api-sdk/stop_event_stream"},
	"stim-status":      {path: "/api-sdk/get_stim_status"},
	"photo-ids":        {path: "/api-sdk/get_image_ids"},
	"photo":            {path: "/api-sdk/get_image", fields: map[string]func(string) bool{"id": integerBetween(0, 2147483647)}},
	"photo-thumb":      {path: "/api-sdk/get_image_thumb", fields: map[string]func(string) bool{"id": integerBetween(0, 2147483647)}},
	"photo-delete":     {path: "/api-sdk/delete_image", fields: map[string]func(string) bool{"id": integerBetween(0, 2147483647)}},
}

func robotSettingsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		writeRobotSettingsError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !sameOrigin(r) {
		writeRobotSettingsError(w, http.StatusForbidden, "request origin not allowed")
		return
	}
	actionName := strings.TrimPrefix(r.URL.Path, "/api/robot-settings/")
	action, ok := robotSettingsActions[actionName]
	if !ok {
		writeRobotSettingsError(w, http.StatusNotFound, "unknown robot setting")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, onboardingMaxBody)
	if err := r.ParseForm(); err != nil {
		writeRobotSettingsError(w, http.StatusBadRequest, "invalid request")
		return
	}
	serial := strings.ToLower(strings.TrimSpace(r.Form.Get("serial")))
	if !validRobotSerial(serial) {
		writeRobotSettingsError(w, http.StatusBadRequest, "invalid robot identifier")
		return
	}
	values := url.Values{"serial": []string{serial}}
	for field, validate := range action.fields {
		value := strings.TrimSpace(r.Form.Get(field))
		if !validate(value) {
			writeRobotSettingsError(w, http.StatusBadRequest, fmt.Sprintf("invalid %s", field))
			return
		}
		values.Set(field, value)
	}
	base, err := url.Parse(onboardingBackend)
	if err != nil {
		writeRobotSettingsError(w, http.StatusServiceUnavailable, "robot service is not configured")
		return
	}
	base.Path = action.path
	request, err := http.NewRequestWithContext(r.Context(), http.MethodPost, base.String(), strings.NewReader(values.Encode()))
	if err != nil {
		writeRobotSettingsError(w, http.StatusBadRequest, "could not prepare robot request")
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := onboardingHTTPClient.Do(request)
	if err != nil {
		writeRobotSettingsError(w, http.StatusBadGateway, "robot service is unavailable")
		return
	}
	defer response.Body.Close()
	contentType := response.Header.Get("Content-Type")
	if actionName == "photo" || actionName == "photo-thumb" || actionName == "camera-frame" {
		contentType = "image/jpeg"
	}
	if contentType == "" {
		contentType = "text/plain; charset=utf-8"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(response.StatusCode)
	_, _ = io.Copy(w, io.LimitReader(response.Body, robotSettingsMaxBody))
}

func validRobotSerial(serial string) bool {
	if len(serial) != 8 {
		return false
	}
	for _, c := range serial {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}

func robotCameraStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	serial := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("serial")))
	if !validRobotSerial(serial) {
		http.Error(w, "invalid robot identifier", http.StatusBadRequest)
		return
	}
	base, err := url.Parse(onboardingBackend)
	if err != nil {
		http.Error(w, "robot service is not configured", http.StatusServiceUnavailable)
		return
	}
	base.Path = "/cam-stream"
	base.RawQuery = url.Values{"serial": []string{serial}}.Encode()
	request, err := http.NewRequestWithContext(r.Context(), http.MethodGet, base.String(), nil)
	if err != nil {
		http.Error(w, "could not prepare camera request", http.StatusBadRequest)
		return
	}
	response, err := (&http.Client{}).Do(request)
	if err != nil {
		http.Error(w, "robot camera is unavailable", http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	if contentType := response.Header.Get("Content-Type"); contentType != "" {
		contentType = strings.Replace(contentType, "boundary=--", "boundary=", 1)
		w.Header().Set("Content-Type", contentType)
	}
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(response.StatusCode)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
		_, _ = io.Copy(flushingWriter{writer: w, flusher: flusher}, response.Body)
		return
	}
	_, _ = io.Copy(w, response.Body)
}

type flushingWriter struct {
	writer  io.Writer
	flusher http.Flusher
}

func (w flushingWriter) Write(data []byte) (int, error) {
	n, err := w.writer.Write(data)
	w.flusher.Flush()
	return n, err
}

func writeRobotSettingsError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	http.Error(w, fmt.Sprintf(`{"error":%q}`, message), status)
}
