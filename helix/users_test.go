package helix

import (
	"encoding/json"
	"net/http"
	"testing"
)

// Tests that received data is empty for a bad request
func TestGetUsersEmpty(t *testing.T) {
	client := newMockTwitchClient("test-id", "test-secret", "client", http.StatusBadRequest, `{"error":"Bad Request","status":400,"message":"Must provide an ID, Login or OAuth Token"}`)

	resp, err := client.GetUsers(&GetUsersOpt{
		Login: "kyrotobi",
	})

	if err != nil {
		t.Error(err)
	}

	if len(resp.Data) != 0 {
		t.Error("expected empty data response")
	}
}

// Test that GetUsers decodes the dummy JSON from the mocked HTTPClient correctly
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
				ProfileImageUrl: "https://static-cdn.jtvnw.net/jtv_user_pictures/dallas-profile_image-1a2c906ee2c35f12-300x300.png",
				OfflineImageUrl: "https://static-cdn.jtvnw.net/jtv_user_pictures/dallas-profile_image-1a2c906ee2c35f12-300x300.png",
				ViewCount:       "999999999",
				Email:           "testemail@gmail.com",
			},
		},
	}
	testRespJson, _ := json.Marshal(testResp)
	client := newMockTwitchClient("test-id", "test-secret", "client", http.StatusOK, string(testRespJson))

	// Doesn't matter what we put here since the response is predetermined
	resp, err := client.GetUsers(&GetUsersOpt{
		Login: "dallas",
	})

	if err != nil {
		t.Error(err)
	}
	if len(resp.Data) != 1 {
		t.Error("expected single data value response")
	}
	if resp.Data[0].Login != testResp.Data[0].Login {
		t.Errorf("got %s, expected %s", resp.Data[0].Login, testResp.Data[0].Login)
	}
	if resp.Data[0].ID != testResp.Data[0].ID {
		t.Errorf("got %s, expected %s", resp.Data[0].ID, testResp.Data[0].ID)
	}
	if resp.Data[0].DisplayName != testResp.Data[0].DisplayName {
		t.Errorf("got %s, expected %s", resp.Data[0].DisplayName, testResp.Data[0].DisplayName)
	}
	if resp.Data[0].Type != testResp.Data[0].Type {
		t.Errorf("got %s, expected %s", resp.Data[0].Type, testResp.Data[0].Type)
	}
	if resp.Data[0].BroadcasterType != testResp.Data[0].BroadcasterType {
		t.Errorf("got %s, expected %s", resp.Data[0].BroadcasterType, testResp.Data[0].BroadcasterType)
	}
	if resp.Data[0].Description != testResp.Data[0].Description {
		t.Errorf("got %s, expected %s", resp.Data[0].Description, testResp.Data[0].Description)
	}
	if resp.Data[0].ProfileImageUrl != testResp.Data[0].ProfileImageUrl {
		t.Errorf("got %s, expected %s", resp.Data[0].ProfileImageUrl, testResp.Data[0].ProfileImageUrl)
	}
	if resp.Data[0].OfflineImageUrl != testResp.Data[0].OfflineImageUrl {
		t.Errorf("got %s, expected %s", resp.Data[0].OfflineImageUrl, testResp.Data[0].OfflineImageUrl)
	}
	if resp.Data[0].ViewCount != testResp.Data[0].ViewCount {
		t.Errorf("got %s, expected %s", resp.Data[0].ViewCount, testResp.Data[0].ViewCount)
	}
	if resp.Data[0].Email != testResp.Data[0].Email {
		t.Errorf("got %s, expected %s", resp.Data[0].Email, testResp.Data[0].Email)
	}
}
