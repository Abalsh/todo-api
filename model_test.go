package todo_api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/goals", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected Empty Array. Got %s", body)
	}
}

func TestGetNonExistentGoal(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/goal/9001", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Goal not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Goal not found'. Got '%s'", m["error"])
	}

}
