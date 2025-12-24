package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var (
	loginURL               = "/auth/login"
	sessionCheckCandidates = []string{
		"/auth/me",
		"/me",
		"/protected",
		"/auth/protected",
		"/api/v1/me",
		"/session",
	}
)

func doLogin(t *testing.T, serverURL string, creds map[string]string) (*http.Response, []byte) {
	t.Helper()
	b, err := json.Marshal(creds)
	if err != nil {
		t.Fatalf("marshal creds: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, serverURL+loginURL, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do login request: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp, body
}

func TestLoginSetsSessionCookie(t *testing.T) {
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	creds := map[string]string{
		"username": "test",
		"password": "test",
	}

	resp, _ := doLogin(t, ts.URL, creds)
	defer resp.Body.Close()

	if len(resp.Cookies()) == 0 {
		t.Fatalf("expected at least one Set-Cookie header on login, got none; status: %d", resp.StatusCode)
	}

	foundNonEmpty := false
	foundHttpOnly := false
	for _, c := range resp.Cookies() {
		if c.Value != "" {
			foundNonEmpty = true
		}
		if c.HttpOnly {
			foundHttpOnly = true
		}
	}
	if !foundNonEmpty {
		t.Fatalf("all cookies had empty values")
	}
	if !foundHttpOnly {
		t.Logf("warning: no HttpOnly cookie found; consider verifying cookie flags in your auth implementation")
	}
}

func TestSessionPersistsAfterLogin(t *testing.T) {
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	creds := map[string]string{
		"username": "test",
		"password": "test",
	}

	resp, _ := doLogin(t, ts.URL, creds)
	defer resp.Body.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
	}
	u := ts.URL
	jar.SetCookies(parseURL(t, u), resp.Cookies())

	var lastErr error
	got200 := false
	for _, ep := range sessionCheckCandidates {
		req, err := http.NewRequest(http.MethodGet, ts.URL+ep, nil)
		if err != nil {
			lastErr = err
			continue
		}
		r, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		r.Body.Close()
		if r.StatusCode == http.StatusOK {
			got200 = true
			break
		}
	}

	if !got200 {
		if lastErr != nil {
			t.Fatalf("session check failed; last error: %v", lastErr)
		}
		t.Fatalf("session check failed; no candidate endpoint returned 200. Tried: %v", sessionCheckCandidates)
	}
}

func parseURL(t *testing.T, s string) *url.URL {
	t.Helper()
	u, err := url.ParseRequestURI(s)
	if err != nil {
		t.Fatalf("parse url %q: %v", s, err)
	}
	return u
}
