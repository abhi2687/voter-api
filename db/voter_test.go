package db_test

import (
	"testing"

	"github.com/abhi2687/voter-api/db"
	"github.com/abhi2687/voter-api/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	voterList, err := db.New()
	assert.Nil(t, err)
	assert.NotNil(t, voterList)
	assert.Equal(t, 0, len(voterList.Voters))
}

func TestAddVoter(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test adding a new voter
	err := voterList.AddVoter(voter1)
	assert.Nil(t, err)
	assert.Equal(t, voter1, voterList.Voters[voter1.VoterId])

	// Test adding another new voter
	err = voterList.AddVoter(voter2)
	assert.Nil(t, err)
	assert.Equal(t, voter2, voterList.Voters[voter2.VoterId])

	// Test adding an existing voter
	err = voterList.AddVoter(voter1)
	assert.NotNil(t, err)
	assert.Equal(t, "voter already exists", err.Error())
}

func TestGetVoter(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test getting a non-existent voter
	_, err := voterList.GetVoter(voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test getting an existing voter
	voterList.AddVoter(voter1)
	voterList.AddVoter(voter2)

	voter, err := voterList.GetVoter(voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, voter1, voter)
}

func TestGetAllVoters(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test getting all voters
	voterList.AddVoter(voter1)
	voterList.AddVoter(voter2)

	voters := voterList.GetAllVoters()
	assert.Equal(t, 2, len(voters))
	assert.Contains(t, voters, voter1)
	assert.Contains(t, voters, voter2)
}

func TestDeleteAllVoters(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test deleting all voters
	voterList.AddVoter(voter1)
	voterList.AddVoter(voter2)

	voterList.DeleteAllVoters()
	assert.Equal(t, 0, len(voterList.Voters))
}

func TestUpdateVoter(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test updating a non-existent voter
	err := voterList.UpdateVoter(voter1, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test updating an existing voter
	voterList.AddVoter(voter1)
	err = voterList.UpdateVoter(voter2, voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, voter2.Name, voterList.Voters[voter1.VoterId].Name)
	assert.Equal(t, voter2.Email, voterList.Voters[voter1.VoterId].Email)
}

func TestDeleteVoter(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	voter2 := testutils.NewRandVoter(2)

	// Test deleting a non-existent voter
	err := voterList.DeleteVoter(voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test deleting an existing voter
	voterList.AddVoter(voter1)
	voterList.AddVoter(voter2)
	err = voterList.DeleteVoter(voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(voterList.Voters))
	assert.NotContains(t, voterList.Voters, voter1.VoterId)
}

func TestAddVoterPoll(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)

	// Test adding a poll to a non-existent voter
	err := voterList.AddVoterPoll(poll1, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test adding a poll to an existing voter
	voterList.AddVoter(voter1)
	err = voterList.AddVoterPoll(poll1, voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(voterList.Voters[voter1.VoterId].VoteHistory))
	assert.Contains(t, voterList.Voters[voter1.VoterId].VoteHistory, poll1)

	// Test adding same poll to an existing voter with existing polls
	err = voterList.AddVoterPoll(poll1, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "poll already exists", err.Error())
}

func TestGetVoterPoll(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)

	// Test getting a poll for a non-existent voter
	_, err := voterList.GetVoterPoll(poll1.PollId, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test getting a poll for an existing voter with no polls
	voterList.AddVoter(voter1)
	_, err = voterList.GetVoterPoll(poll1.PollId, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "poll does not exist", err.Error())

	// test getting a poll for an existing voter with polls
	voterList.AddVoterPoll(poll1, voter1.VoterId)

	voterPoll, err := voterList.GetVoterPoll(poll1.PollId, voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, poll1, voterPoll)
}

func TestGetVoterPolls(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)

	// Test getting polls for a non-existent voter
	_, err := voterList.GetVoterPolls(voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test getting polls for an existing voter
	voterList.AddVoter(voter1)
	voterList.AddVoterPoll(poll1, voter1.VoterId)

	voterPolls, err := voterList.GetVoterPolls(voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(voterPolls))
	assert.Contains(t, voterPolls, poll1)
}

func TestDeleteVoterPoll(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)

	// Test deleting polls for a non-existent voter
	err := voterList.DeleteVoterPoll(poll1.PollId, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test deleting polls for an existing voter with no polls
	voterList.AddVoter(voter1)
	err = voterList.DeleteVoterPoll(poll1.PollId, voter1.VoterId)
	assert.NotNil(t, err)
	assert.Equal(t, "poll does not exist", err.Error())

	// test deleting polls for an existing voter with polls
	voterList.AddVoterPoll(poll1, voter1.VoterId)
	err = voterList.DeleteVoterPoll(poll1.PollId, voter1.VoterId)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(voterList.Voters[voter1.VoterId].VoteHistory))
}

func TestUpdateVoterPoll(t *testing.T) {
	voterList := &db.VoterList{
		Voters: make(map[uint]db.Voter),
	}

	voter1 := testutils.NewRandVoter(1)
	poll1 := testutils.NewRandPollVoteRecord(1)
	poll2 := testutils.NewRandPollVoteRecord(2)

	// Test updating polls for a non-existent voter
	err := voterList.UpdateVoterPoll(poll1, voter1.VoterId, poll1.PollId)
	assert.NotNil(t, err)
	assert.Equal(t, "voter does not exist", err.Error())

	// Test updating polls for an existing voter with no polls
	voterList.AddVoter(voter1)
	err = voterList.UpdateVoterPoll(poll1, voter1.VoterId, poll1.PollId)
	assert.NotNil(t, err)
	assert.Equal(t, "poll does not exist", err.Error())

	// test updating polls for an existing voter with polls
	voterList.AddVoterPoll(poll1, voter1.VoterId)
	err = voterList.UpdateVoterPoll(poll2, voter1.VoterId, poll1.PollId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(voterList.Voters[voter1.VoterId].VoteHistory))
	assert.Equal(t, voterList.Voters[voter1.VoterId].VoteHistory[0].VoteId, poll2.VoteId)
	assert.Equal(t, voterList.Voters[voter1.VoterId].VoteHistory[0].VoteDate, poll2.VoteDate)
}
