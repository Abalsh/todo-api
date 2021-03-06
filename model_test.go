package todo_api

import (
	"bytes"
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

func TestCreateGoal(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`{"name":"really bad goal", "description": "don't buy"}`)
	req, _ := http.NewRequest("POST", "/goal", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "really bad goal" {
		t.Errorf("Expected Goal name to be 'really bad goal'. Got '%v'", m["name"])
	}
	if m["description"] != "don't buy" {
		t.Errorf("Expected Goal name to be `don't buy`. Got `%v`", m["description"])
	}
	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetGoal(t *testing.T) {
	clearTable()
	addGoal(1)
	req, _ := http.NewRequest("GET", "/goal/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateGoal(t *testing.T) {
	clearTable()
	addGoal(1)

	req, _ := http.NewRequest("GET", "/goal/1", nil)
	response := executeRequest(req)
	var originalGoal map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalGoal)

	var jsonStr = []byte(`{"name":"test goal updated", "description":"description updated"}`)
	req, _ = http.NewRequest("PUT", "/goal/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalGoal["id"] {
		t.Errorf("Expected the id to remain the same (%v)!!!! Got %v instead!!", originalGoal["id"], m["id"])
	}
	if m["name"] == originalGoal["name"] {
		t.Errorf("Expected the name to change from %v to  %v !! got '%v'", originalGoal["name"], m["name"], m["name"])
	}
	if m["description"] == originalGoal["description"] {
		t.Errorf("Expected the description to change from %v to  %v !! got '%v'", originalGoal["description"], m["description"], m["description"])
	}
}
func TestDeleteGoal(t *testing.T) {
	clearTable()
	addGoal(1)

	req, _ := http.NewRequest("GET", "/goal/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/goal/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/goal/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

}
