package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/abhi2687/voter-api/api"
	"github.com/abhi2687/voter-api/db"
	"github.com/abhi2687/voter-api/testutils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var (
	// BASE_API = "http://localhost:1080"

	// cli         = resty.New()
	app             = fiber.New()
	voterHandler, _ = api.New()
)

func init() {
	app.Post("/voters", voterHandler.AddVoter)
	app.Delete("/voters", voterHandler.DeleteAllVoters)
	app.Get("/voters/:id", voterHandler.GetVoter)
	app.Get("/voters", voterHandler.GetAllVoters)
	app.Put("/voters/:id", voterHandler.UpdateVoter)
	app.Delete("/voters/:id", voterHandler.DeleteVoter)
	app.Get("/voters/:id/polls", voterHandler.GetVoterPolls)
	app.Post("/voters/:id/polls", voterHandler.AddVoterPoll)
	app.Get("/voters/:id/polls/:pollid", voterHandler.GetVoterPoll)
	app.Put("/voters/:id/polls/:pollid", voterHandler.UpdateVoterPoll)
	app.Delete("/voters/:id/polls/:pollid", voterHandler.DeleteVoterPoll)
}

func deleteAllVoters() {
	req, err := http.NewRequest("DELETE", "/voters", nil)
	if err != nil {
		fmt.Printf("failed to create HTTP request to delete all voters: %v", err)
		return
	}

	// Serve the request
	_, err = app.Test(req)
	if err != nil {
		fmt.Printf("failed to delete all voters: %v", err)
	}
}

// testing voter handler New function
func TestNew(t *testing.T) {
	voterHandler, err := api.New()
	if err != nil {
		t.Errorf("error creating voter handler: %v", err)
	}
	if voterHandler == nil {
		t.Errorf("error creating voter handler: %v", err)
	}
}

// testing voter handler AddVoter - Invalid body
func TestAddVoterBadRequest(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}

	// Serve the request - without content type will result in error
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// testing voter handler AddVoter - Voter exists
func TestAddVoterVoterExists(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("failed to serve request: %v", err)
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// testing voter handler AddVoter - Success case
func TestAddVoterSuccessful(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Unmarshal the response body
	var responseVoter db.Voter
	err = json.Unmarshal(body, &responseVoter)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response body
	assert.Equal(t, voter, responseVoter)
}

// testing voter handler GetVoter -  Success case
func TestGetVoter(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request - Get voter and verify
	req, err = http.NewRequest("GET", fmt.Sprintf("/voters/%d", voter.VoterId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body
	var responseVoter db.Voter
	err = json.Unmarshal(body, &responseVoter)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response body
	assert.Equal(t, voter, responseVoter)
}

// testing voter handler GetAllVoters - Success case
func TestGetAllVoter(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Marshal the Voter to JSON
	voter1JSON, err := json.Marshal(voter1)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	voter2JSON, err := json.Marshal(voter2)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voter1JSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter1
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("POST", "/voters", bytes.NewBuffer(voter2JSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter2
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request - Get all voters and verify
	req, err = http.NewRequest("GET", "/voters", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body
	var voters []db.Voter
	err = json.Unmarshal(body, &voters)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response body
	assert.Contains(t, voters, voter1)
	assert.Contains(t, voters, voter2)
}

// testing voter handler UpdateVoter - Success case
func TestUpdateVoter(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	updateVoter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	updateVoterJSON, err := json.Marshal(updateVoter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("PUT", fmt.Sprintf("/voters/%d", voter.VoterId), bytes.NewBuffer(updateVoterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Update voter
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Create a new HTTP request - Get voter and verify
	req, err = http.NewRequest("GET", fmt.Sprintf("/voters/%d", voter.VoterId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body
	var responseVoter db.Voter
	err = json.Unmarshal(body, &responseVoter)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response body
	assert.Equal(t, updateVoter, responseVoter)
}

// testing voter handler UpdateVoter - Invalid body
func TestUpdateVoterBadRequest(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	updateVoter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	updateVoterJSON, err := json.Marshal(updateVoter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("PUT", fmt.Sprintf("/voters/%d", voter.VoterId), bytes.NewBuffer(updateVoterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Update voter without setting content type
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// testing voter handler UpdateVoter - Updating voter that doesnt exists
func TestUpdateVoterNotExists(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	updateVoterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("PUT", fmt.Sprintf("/voters/%d", voter.VoterId), bytes.NewBuffer(updateVoterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Update voter
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// testing voter handler DeleteVoter - Success case
func TestDeleteVoter(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/voters/%d", voter.VoterId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Delete voter
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Create a new HTTP request - Get voter and verify
	req, err = http.NewRequest("GET", fmt.Sprintf("/voters/%d", voter.VoterId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// testing voter handler DeleteVoter - Voter doesnt exists
func TestDeleteVoterNotExists(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/voters/%d", 1), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Delete voter
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// testing voter handler DeleteAllVoters function by mocking the API call
func TestDeleteAllVoters(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Marshal the Voter to JSON
	voter1JSON, err := json.Marshal(voter1)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	voter2JSON, err := json.Marshal(voter2)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voter1JSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter1
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("POST", "/voters", bytes.NewBuffer(voter2JSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter2
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("DELETE", "/voters", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Delete all voters
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Create a new HTTP request - Get all voters and verify
	req, err = http.NewRequest("GET", "/voters", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body
	var voters []db.Voter
	err = json.Unmarshal(body, &voters)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response body
	assert.Empty(t, voters)
}

// testing voter handler GetVoterPolls - Success case
func TestGetVoterPolls(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)
	poll := testutils.NewRandPollVoteRecord(1)
	voter.VoteHistory = append(voter.VoteHistory, poll)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("GET", fmt.Sprintf("/voters/%d/polls", voter.VoterId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Get voter polls
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body
	var voterPolls []db.VoterHistory
	err = json.Unmarshal(body, &voterPolls)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response body
	assert.Contains(t, voterPolls, poll)
}

// testing voter handler GetVoterPolls - Voter doesnt exists
func TestGetVoterPollsNotExists(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fmt.Sprintf("/voters/%d/polls", 1), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Get voter polls
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// testing voter handler AddVoterPoll - Success case
func TestAddVoterPoll(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)
	poll := testutils.NewRandPollVoteRecord(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Marshal the Poll to JSON
	pollJSON, err := json.Marshal(poll)
	if err != nil {
		t.Fatalf("Failed to marshal poll to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("POST", fmt.Sprintf("/voters/%d/polls", voter.VoterId), bytes.NewBuffer(pollJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter poll
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Create a new HTTP request - Get voter polls and verify
	req, err = http.NewRequest("GET", fmt.Sprintf("/voters/%d/polls", voter.VoterId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body
	var voterPolls []db.VoterHistory
	err = json.Unmarshal(body, &voterPolls)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response body
	assert.Contains(t, voterPolls, poll)
}

// testing voter handler AddVoterPoll - Invalid Body
func TestAddVoterPollBadRequest(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)
	poll := testutils.NewRandPollVoteRecord(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Marshal the Poll to JSON
	pollJSON, err := json.Marshal(poll)
	if err != nil {
		t.Fatalf("Failed to marshal poll to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("POST", fmt.Sprintf("/voters/%d/polls", voter.VoterId), bytes.NewBuffer(pollJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Add voter poll without content type header
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// testing voter handler AddVoterPoll - Poll already exists
func TestAddVoterPollExists(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)
	poll := testutils.NewRandPollVoteRecord(1)
	voter.VoteHistory = append(voter.VoteHistory, poll)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Marshal the Poll to JSON
	pollJSON, err := json.Marshal(poll)
	if err != nil {
		t.Fatalf("Failed to marshal poll to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("POST", fmt.Sprintf("/voters/%d/polls", voter.VoterId), bytes.NewBuffer(pollJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter poll
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetVoterPoll(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)
	poll := testutils.NewRandPollVoteRecord(1)
	voter.VoteHistory = append(voter.VoteHistory, poll)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request - Get voter poll and verify
	req, err = http.NewRequest("GET", fmt.Sprintf("/voters/%d/polls/%d", voter.VoterId, poll.PollId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body
	var voterPoll db.VoterHistory
	err = json.Unmarshal(body, &voterPoll)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response body
	assert.Equal(t, poll, voterPoll)
}

// testing voter handler UpdateVoterPoll - Success case
func TestUpdateVoterPoll(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)
	poll := testutils.NewRandPollVoteRecord(1)
	voter.VoteHistory = append(voter.VoteHistory, poll)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	updatePoll := testutils.NewRandPollVoteRecord(1)

	// Marshal the Poll to JSON
	updatePollJSON, err := json.Marshal(updatePoll)
	if err != nil {
		t.Fatalf("Failed to marshal poll to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("PUT", fmt.Sprintf("/voters/%d/polls/%d", voter.VoterId, poll.PollId), bytes.NewBuffer(updatePollJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Update voter poll
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Create a new HTTP request - Get voter poll and verify
	req, err = http.NewRequest("GET", fmt.Sprintf("/voters/%d/polls/%d", voter.VoterId, poll.PollId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body
	var responsePoll db.VoterHistory
	err = json.Unmarshal(body, &responsePoll)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response body
	assert.Equal(t, updatePoll, responsePoll)
}

// testing voter handler UpdateVoterPoll - BadRequest
func TestUpdateVoterPollBadRequest(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)
	poll := testutils.NewRandPollVoteRecord(1)
	voter.VoteHistory = append(voter.VoteHistory, poll)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	updatePoll := testutils.NewRandPollVoteRecord(1)

	// Marshal the Poll to JSON
	updatePollJSON, err := json.Marshal(updatePoll)
	if err != nil {
		t.Fatalf("Failed to marshal poll to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("PUT", fmt.Sprintf("/voters/%d/polls/%d", voter.VoterId, poll.PollId), bytes.NewBuffer(updatePollJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Update voter poll without content type header
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// testing voter handler UpdateVoterPoll - Poll doesnt exists
func TestUpdateVoterPollNotExists(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	poll := testutils.NewRandPollVoteRecord(1)

	// Marshal the Poll to JSON
	updatePollJSON, err := json.Marshal(poll)
	if err != nil {
		t.Fatalf("Failed to marshal poll to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("PUT", fmt.Sprintf("/voters/%d/polls/%d", voter.VoterId, poll.PollId), bytes.NewBuffer(updatePollJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Update voter poll
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// testing voter handler DeleteVoterPoll - Success case
func TestDeleteVoterPoll(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)
	poll := testutils.NewRandPollVoteRecord(1)
	voter.VoteHistory = append(voter.VoteHistory, poll)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/voters/%d/polls/%d", voter.VoterId, poll.PollId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Delete voter poll
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Create a new HTTP request - Get voter poll and verify
	req, err = http.NewRequest("GET", fmt.Sprintf("/voters/%d/polls/%d", voter.VoterId, poll.PollId), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// testing voter handler DeleteVoterPoll - Poll doesnt exists
func TestDeleteVoterPollNotExists(t *testing.T) {
	// clean up existing voters
	deleteAllVoters()

	// Create a new Voter
	voter := testutils.NewRandVoter(1)

	// Marshal the Voter to JSON
	voterJSON, err := json.Marshal(voter)
	if err != nil {
		t.Fatalf("Failed to marshal voter to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/voters", bytes.NewBuffer(voterJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	// Serve the request - Add voter
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Create a new HTTP request
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/voters/%d/polls/%d", voter.VoterId, 1), nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Serve the request - Delete voter poll
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to serve request: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// func TestHealthCheck(t *testing.T) {
// 	// Create a new HTTP request
// 	req, err := http.NewRequest("GET", "/health", nil)
// 	if err != nil {
// 		t.Fatalf("Failed to create HTTP request: %v", err)
// 	}

// 	// Serve the request
// 	resp, err := app.Test(req)
// 	if err != nil {
// 		t.Fatalf("Failed to serve request: %v", err)
// 	}

// 	// Check the status code
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)
// }

//

// func TestGetAllVoter(t *testing.T) {
// 	var voters []db.Voter

// 	// clean up existing voters
// 	_, err := cli.R().Delete(BASE_API + "/voters")
// 	if err != nil {
// 		log.Printf("error deleting all voters, %v", err)
// 		os.Exit(1)
// 	}

// 	voter1 := newRandVoter(1)
// 	voter2 := newRandVoter(2)

// 	_, err1 := cli.R().SetBody(voter1).Post(BASE_API + "/voters")
// 	if err1 != nil {
// 		log.Printf("error adding voter1, %v", err1)
// 		os.Exit(1)
// 	}

// 	_, err2 := cli.R().SetBody(voter2).Post(BASE_API + "/voters")
// 	if err2 != nil {
// 		log.Printf("error adding voter1, %v", err2)
// 		os.Exit(1)
// 	}

// 	rsp, err := cli.R().SetResult(&voters).Get(BASE_API + "/voters")
// 	if err != nil {
// 		log.Printf("error getting all voters, %v", err)
// 		os.Exit(1)
// 	}

// 	if rsp.StatusCode() != 200 {
// 		log.Printf("error getting all voters, %v", err)
// 		os.Exit(1)
// 	}

// 	assert.Nil(t, err)
// 	assert.Equal(t, 200, rsp.StatusCode())
// 	assert.Equal(t, "application/json", rsp.Header().Get("Content-Type"))
// 	assert.Equal(t, 2, len(voters))
// 	assert.Contains(t, voters, voter1)
// 	assert.Contains(t, voters, voter2)
// }
