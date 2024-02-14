package db

import (
	"errors"
	"time"
)

type VoterHistory struct {
	PollId   uint      `json:"pollId"`
	VoteId   uint      `json:"voteId"`
	VoteDate time.Time `json:"voteDate"`
}

type Voter struct {
	VoterId     uint           `json:"voterId"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	VoteHistory []VoterHistory `json:"voteHistory,omitempty"`
}

type VoterList struct {
	Voters map[uint]Voter //A map of VoterIDs as keys and Voter structs as values
}

func New() (*VoterList, error) {
	return &VoterList{
		Voters: make(map[uint]Voter),
	}, nil
}

func (v *VoterList) AddVoter(voter Voter) error {
	_, ok := v.Voters[voter.VoterId]
	if ok {
		return errors.New("voter already exists")
	}

	v.Voters[voter.VoterId] = voter
	return nil
}

func (v *VoterList) GetVoter(voterId uint) (Voter, error) {
	voter, ok := v.Voters[voterId]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}

	return voter, nil
}

func (v *VoterList) GetAllVoters() []Voter {
	var voterList []Voter
	for _, voter := range v.Voters {
		voterList = append(voterList, voter)
	}
	return voterList
}

func (v *VoterList) DeleteAllVoters() {
	v.Voters = make(map[uint]Voter)
}

func (v *VoterList) UpdateVoter(voter Voter, voterId uint) error {
	_, ok := v.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	updatedVoter := v.Voters[voterId]
	updatedVoter.Email = voter.Email
	updatedVoter.Name = voter.Name

	v.Voters[voterId] = updatedVoter

	return nil
}

func (v *VoterList) DeleteVoter(voterId uint) error {
	_, ok := v.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	delete(v.Voters, voterId)
	return nil
}

func (v *VoterList) GetVoterPolls(voterId uint) ([]VoterHistory, error) {
	voter, ok := v.Voters[voterId]
	if !ok {
		return nil, errors.New("voter does not exist")
	}

	return voter.VoteHistory, nil
}

func (v *VoterList) AddVoterPoll(voterPoll VoterHistory, voterId uint) error {
	voter, ok := v.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	for _, vh := range voter.VoteHistory {
		if vh.PollId == voterPoll.PollId {
			return errors.New("poll already exists")
		}
	}

	voter.VoteHistory = append(voter.VoteHistory, voterPoll)

	v.Voters[voterId] = voter
	return nil
}

func (v *VoterList) GetVoterPoll(voterId uint, pollId uint) (VoterHistory, error) {
	voter, ok := v.Voters[voterId]
	if !ok {
		return VoterHistory{}, errors.New("voter does not exist")
	}

	for _, vh := range voter.VoteHistory {
		if vh.PollId == pollId {
			return vh, nil
		}
	}

	return VoterHistory{}, errors.New("poll does not exist")
}

func (v *VoterList) UpdateVoterPoll(voterPoll VoterHistory, voterId uint, pollId uint) error {
	voter, ok := v.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	for i, vh := range voter.VoteHistory {
		if vh.PollId == pollId {
			voterPoll.PollId = pollId
			voter.VoteHistory[i] = voterPoll
			return nil
		}
	}

	return errors.New("poll does not exist")
}

func (v *VoterList) DeleteVoterPoll(voterId uint, pollId uint) error {
	voter, ok := v.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	for i, vh := range voter.VoteHistory {
		if vh.PollId == pollId {
			voter.VoteHistory = append(voter.VoteHistory[:i], voter.VoteHistory[i+1:]...)
			v.Voters[voterId] = voter
			return nil
		}
	}

	return errors.New("poll does not exist")
}
