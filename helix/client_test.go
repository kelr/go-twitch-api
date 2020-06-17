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
func newMockResponse(status int, body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(body)
	}
}

// Create a mocked Client with a mocked HTTPClient
func newMockClient(clientID string, clientSecret string, tokenType string, respStatus int, respBody []byte) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		conn: &mockHTTPClient{
			response: newMockResponse(respStatus, respBody),
		},
		tokenType: tokenType,
	}
}
