package helix

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Tests that received data is empty for a bad request
func TestGetUsersEmpty(t *testing.T) {
	cfg := new(Config)
	client := newMockClient(cfg, "client", http.StatusBadRequest, []byte(`{"error":"Bad Request","status":400,"message":"Must provide an ID, Login or OAuth Token"}`))

	resp, err := client.GetUsers(&GetUsersOpt{
		Login: []string{"kyrotobi"},
	})

	if err != nil {
		t.Error(err)
	}

	if len(resp.Data) != 0 {
		t.Error("expected empty data response")
	}
}

// Test that GetUsers decodes the JSON from the internal HTTPClient correctly
func TestGetUsers(t *testing.T) {
	testResp := &GetUsersResponse{
		Data: []GetUsersData{
			{
				ID:              "123123",
				Login:           "testlogin",
				DisplayName:     "testdisplayname",
				Type:            "",
				BroadcasterType: "partner",
				Description:     "hi im strimmer :)",
				ProfileImageURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/dallas-profile_image-1a2c906ee2c35f12-300x300.png",
				OfflineImageURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/dallas-profile_image-1a2c906ee2c35f12-300x300.png",
				ViewCount:       999999999,
				Email:           "testemail@gmail.com",
			},
		},
	}
	testRespJSON, _ := json.Marshal(testResp)
	client := newMockClient(new(Config), "app", http.StatusOK, testRespJSON)

	// Doesn't matter what we put here.
	resp, err := client.GetUsers(&GetUsersOpt{
		Login: []string{"dallas"},
	})

	if err != nil {
		t.Error(err)
	}
	if len(resp.Data) != 1 {
		t.Error("expected single data value response")
	}
	if !cmp.Equal(*testResp, *resp) {
		t.Error("decoded struct does not match input struct")
	}
}
