package testutils

import (
	"github.com/abhi2687/voter-api/db"
	fake "github.com/brianvoe/gofakeit/v6"
)

func NewRandVoter(id uint) db.VoterItem {
	return db.VoterItem{
		VoterId: id,
		Name:    fake.Name(),
		Email:   fake.Email(),
	}
}

func NewRandPollVoteRecord(pollId uint) db.VoterHistory {
	return db.VoterHistory{
		PollId:   pollId,
		VoteId:   fake.UintRange(1, 100),
		VoteDate: fake.Date(),
	}
}
