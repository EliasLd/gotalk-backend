package http

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/google/uuid"
)


func ParseUserIDFromResponse(t *testing.T, body io.Reader) uuid.UUID {
	t.Helper()

	var response map[string]interface{}
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	idRaw, ok := response["id"]
	if !ok {
		t.Fatalf("Response does not contain 'id' field")
	}

	idStr, ok := idRaw.(string)
	if !ok {
		t.Fatalf("Expected 'id' field to be a string, got %T", idRaw)
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		t.Fatalf("Invalid UUID format: %v", err)
	}

	return userID
}
