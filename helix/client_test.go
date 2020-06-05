package helix

import (
	"net/http"
	"net/http/httptest"
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
func newMockResponse(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(body))
	}
}

// Create a mocked TwitchClient with a mocked HTTPClient
func newMockTwitchClient(clientId string, clientSecret string, tokenType string, respStatus int, respBody string) *TwitchClient {
	return &TwitchClient{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		conn: &mockHTTPClient{
			response: newMockResponse(respStatus, respBody),
		},
		tokenType: tokenType,
	}
}
