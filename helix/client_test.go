package helix

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// A fake HTTPClient that implements the HTTPClient interface for testing
type mockHTTPClient struct {
	response http.HandlerFunc
}

// Implement the Do method to satisfy the HTTPClient interface
func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(m.response)
	handler.ServeHTTP(rec, req)
	return rec.Result(), nil
}

// Create a HandlerFunc with the HTTP status code and body that the mock HTTP Client will respond with
func newMockResponse(status int, body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(body)
	}
}

// Create a mocked Client with a mocked HTTPClient
func newMockClient(cfg *Config, tokenType string, respStatus int, respBody []byte) *Client {
	return &Client{
		conn: &mockHTTPClient{
			response: newMockResponse(respStatus, respBody),
		},
		config:    cfg,
		tokenType: tokenType,
	}
}

// Tests that the HTTP request and header are properly constructed.
func TestBuildRequest(t *testing.T) {
	cases := []struct {
		inputOpts        interface{}
		inputBaseURL     string
		expectedClientID string
		expectedURL      string
		expectedMethod   string
	}{
		{
			inputOpts: &GetUsersOpt{
				ID:    []string{"123456789", "987654321"},
				Login: []string{"dallas", "kyrotobi"},
			},
			inputBaseURL:     getUsersPath,
			expectedClientID: "test-client-id",
			expectedURL:      "https://api.twitch.tv/helix/users?id=123456789&id=987654321&login=dallas&login=kyrotobi",
			expectedMethod:   http.MethodGet,
		},
		{
			inputOpts:        nil,
			inputBaseURL:     getUsersPath,
			expectedClientID: "test-client-id",
			expectedURL:      "https://api.twitch.tv/helix/users",
			expectedMethod:   http.MethodGet,
		},
		{
			inputOpts: &GetUsersOpt{
				Login: []string{"kyrotobi", "kyrotobi", "kyrotobi", "kyrotobi", "kyrotobi", "kyrotobi"},
			},
			inputBaseURL:     getUsersPath,
			expectedClientID: "",
			expectedURL:      "https://api.twitch.tv/helix/users?login=kyrotobi&login=kyrotobi&login=kyrotobi&login=kyrotobi&login=kyrotobi&login=kyrotobi",
			expectedMethod:   http.MethodGet,
		},
	}

	for _, c := range cases {
		client := newMockClient(&Config{
			ClientID: c.expectedClientID,
		}, "app", http.StatusOK, nil)
		got, err := client.buildRequest(c.inputBaseURL, c.inputOpts, c.expectedMethod)
		if err != nil {
			t.Error(err)
		}
		if c.expectedClientID != got.Header["Client-Id"][0] {
			t.Errorf("wanted: %s\n got: %s\n", c.expectedClientID, got.Header["Client-Id"][0])
		}
		if c.expectedURL != got.URL.String() {
			t.Errorf("wanted: %s\n got: %s\n", c.expectedURL, got.URL.String())
		}
		if c.expectedMethod != got.Method {
			t.Errorf("wanted: %s\n got: %s\n", c.expectedMethod, got.Method)
		}
	}
}
