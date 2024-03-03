package db_test

import (
	"testing"

	"github.com/abhi2687/voter-api/db"
	"github.com/abhi2687/voter-api/testutils"
	"github.com/stretchr/testify/assert"
)

// Initialize a new voter
var voter, err = db.New()

func TestNew(t *testing.T) {
	assert.Nil(t, err)
	assert.NotNil(t, voter)
}

func TestAddVoter(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test adding a new voter
	err := voter.AddVoter(voter1)
	assert.Nil(t, err)

	voterFromRedis, _ := voter.GetVoter(voter1.VoterId)
	assert.Equal(t, voter1, voterFromRedis)

	// Test adding another new voter
	err = voter.AddVoter(voter2)
	assert.Nil(t, err)
	voterFromRedis, _ = voter.GetVoter(voter2.VoterId)
	assert.Equal(t, voter2, voterFromRedis)

	// Test adding an existing voter
	err = voter.AddVoter(voter1)
	assert.NotNil(t, err)
	assert.Equal(t, "voter already exists", err.Error())
}

func TestGetVoter(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test getting a non-existent voter
	voterItem, err := voter.GetVoter(voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, db.VoterItem{}, voterItem)

	// Test getting an existing voter
	voter.AddVoter(voter1)
	voter.AddVoter(voter2)

	voterFromRedis, err := voter.GetVoter(voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, voter1, voterFromRedis)
}

func TestGetAllVoters(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test getting all voters
	voter.AddVoter(voter1)
	voter.AddVoter(voter2)

	voters := voter.GetAllVoters()
	assert.Equal(t, 2, len(voters))
	assert.Contains(t, voters, voter1)
	assert.Contains(t, voters, voter2)
}

func TestDeleteAllVoters(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test deleting all voters
	voter.AddVoter(voter1)
	voter.AddVoter(voter2)

	voter.DeleteAllVoters()
	voters := voter.GetAllVoters()
	assert.Equal(t, 0, len(voters))
}

func TestUpdateVoter(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Updating voter2 to have the same voterId as voter1
	voter2.VoterId = voter1.VoterId

	// Test updating a non-existent voter
	err := voter.UpdateVoter(voter1, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test updating an existing voter
	voter.AddVoter(voter1)
	err = voter.UpdateVoter(voter2, voter1.VoterId)
	assert.Nil(t, err)

	// Get the updated voter
	voter, _ := voter.GetVoter(voter1.VoterId)
	assert.Equal(t, voter2.Name, voter.Name)
	assert.Equal(t, voter2.Email, voter.Email)
}

func TestDeleteVoter(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test deleting a non-existent voter
	err := voter.DeleteVoter(voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test deleting an existing voter
	voter.AddVoter(voter1)
	voter.AddVoter(voter2)
	err = voter.DeleteVoter(voter1.VoterId)
	assert.Nil(t, err)

	// Get all voters
	voters := voter.GetAllVoters()

	assert.Equal(t, 1, len(voters))
	assert.NotContains(t, voters, voter1.VoterId)
}

func TestAddVoterPoll(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)

	// Test adding a poll to a non-existent voter
	err := voter.AddVoterPoll(poll1, voter1.VoterId)
	assert.NotNil(t, err)

	// Test adding a poll to an existing voter
	voter.AddVoter(voter1)
	err = voter.AddVoterPoll(poll1, voter1.VoterId)
	assert.Nil(t, err)

	// Test adding same poll to an existing voter with existing polls
	err = voter.AddVoterPoll(poll1, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "poll already exists", err.Error())

	// Get the voter
	voter, _ := voter.GetVoter(voter1.VoterId)
	assert.Equal(t, 1, len(voter.VoteHistory))
	assert.Contains(t, voter.VoteHistory, poll1)
}

func TestGetVoterPoll(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)

	// Test getting a poll for a non-existent voter
	_, err := voter.GetVoterPoll(poll1.PollId, voter1.VoterId)
	assert.NotNil(t, err)

	// Test getting a poll for an existing voter with no polls
	voter.AddVoter(voter1)
	_, err = voter.GetVoterPoll(poll1.PollId, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "poll does not exist", err.Error())

	// test getting a poll for an existing voter with polls
	voter.AddVoterPoll(poll1, voter1.VoterId)

	voterPoll, err := voter.GetVoterPoll(poll1.PollId, voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, poll1, voterPoll)
}

func TestGetVoterPolls(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)

	// Test getting polls for a non-existent voter
	_, err := voter.GetVoterPolls(voter1.VoterId)
	assert.NotNil(t, err)

	// Test getting polls for an existing voter
	voter.AddVoter(voter1)
	voter.AddVoterPoll(poll1, voter1.VoterId)

	voterPolls, err := voter.GetVoterPolls(voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(voterPolls))
	assert.Contains(t, voterPolls, poll1)
}

func TestDeleteVoterPoll(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)

	// Test deleting polls for a non-existent voter
	err := voter.DeleteVoterPoll(poll1.PollId, voter1.VoterId)
	assert.NotNil(t, err)

	// Test deleting polls for an existing voter with no polls
	voter.AddVoter(voter1)
	err = voter.DeleteVoterPoll(poll1.PollId, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "poll does not exist", err.Error())

	// test deleting polls for an existing voter with polls
	voter.AddVoterPoll(poll1, voter1.VoterId)
	err = voter.DeleteVoterPoll(poll1.PollId, voter1.VoterId)
	assert.Nil(t, err)

	// Get the voter
	voter, _ := voter.GetVoter(voter1.VoterId)
	assert.Equal(t, 0, len(voter.VoteHistory))
	assert.NotContains(t, voter.VoteHistory, poll1)
}

func TestUpdateVoterPoll(t *testing.T) {
	//Delete all voters
	voter.DeleteAllVoters()

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)
	poll2 := testutils.NewRandPollVoteRecord(2)

	// Test updating polls for a non-existent voter
	err := voter.UpdateVoterPoll(poll1, voter1.VoterId, poll1.PollId)
	assert.NotNil(t, err)

	// Test updating polls for an existing voter with no polls
	voter.AddVoter(voter1)
	err = voter.UpdateVoterPoll(poll1, voter1.VoterId, poll1.PollId)
	assert.NotNil(t, err)
	assert.Equal(t, "poll does not exist", err.Error())

	// test updating polls for an existing voter with polls
	voter.AddVoterPoll(poll1, voter1.VoterId)
	err = voter.UpdateVoterPoll(poll2, voter1.VoterId, poll1.PollId)
	assert.Nil(t, err)

	// Get the voter
	voter, _ := voter.GetVoter(voter1.VoterId)
	assert.Equal(t, 1, len(voter.VoteHistory))
	assert.Contains(t, voter.VoteHistory, poll2)
}
