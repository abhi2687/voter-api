package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
)

type VoterHistory struct {
	PollId   uint      `json:"pollId"`
	VoteId   uint      `json:"voteId"`
	VoteDate time.Time `json:"voteDate"`
}

type VoterItem struct {
	VoterId     uint           `json:"voterId"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	VoteHistory []VoterHistory `json:"voteHistory,omitempty"`
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type Voter struct {
	//Redis cache connections
	cache
}

func New() (*Voter, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	fmt.Println("DEBUG:  USING REDIS URL: " + redisUrl)
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*Voter, error) {
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		fmt.Println("Error connecting to redis" + err.Error() + "cache might not be available, continuing...")
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	return &Voter{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

func redisKeyFromId(id int) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (v *Voter) getItemFromRedis(key string, voterItem *VoterItem) error {
	itemObject, err := v.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), voterItem)
	if err != nil {
		return err
	}

	return nil
}

func (v *Voter) AddVoter(voterItem VoterItem) error {
	//Check if Voter already exists
	redisKey := redisKeyFromId(int(voterItem.VoterId))

	var existingItem VoterItem
	if err := v.getItemFromRedis(redisKey, &existingItem); err == nil {
		return errors.New("voter already exists")
	}

	//Add item to database with JSON Set
	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voterItem); err != nil {
		return err
	}

	return nil
}

func (v *Voter) GetVoter(voterId uint) (VoterItem, error) {
	redisKey := redisKeyFromId(int(voterId))

	var voterItem VoterItem
	if err := v.getItemFromRedis(redisKey, &voterItem); err != nil {
		return VoterItem{}, err
	}

	return voterItem, nil
}

func (v *Voter) GetAllVoters() []VoterItem {
	keys, err := v.cacheClient.Keys(v.context, RedisKeyPrefix+"*").Result()
	if err != nil {
		fmt.Println("Error getting keys from redis: " + err.Error())
		return nil
	}

	var voterItems []VoterItem
	for _, key := range keys {
		var voterItem VoterItem
		if err := v.getItemFromRedis(key, &voterItem); err != nil {
			fmt.Println("Error getting voters from redis: " + err.Error())
			continue
		}
		voterItems = append(voterItems, voterItem)
	}

	return voterItems
}

func (v *Voter) DeleteAllVoters() {
	keys, err := v.cacheClient.Keys(v.context, RedisKeyPrefix+"*").Result()
	if err != nil {
		fmt.Println("Error getting keys from redis: " + err.Error())
		return
	}

	for _, key := range keys {
		if err := v.cacheClient.Del(v.context, key).Err(); err != nil {
			fmt.Println("Error deleting voters from redis: " + err.Error())
		}
	}
}

func (v *Voter) UpdateVoter(voterItem VoterItem, voterId uint) error {
	redisKey := redisKeyFromId(int(voterItem.VoterId))
	var existingItem VoterItem
	if err := v.getItemFromRedis(redisKey, &existingItem); err != nil {
		return errors.New("voter does not exist")
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voterItem); err != nil {
		return err
	}

	return nil
}

func (v *Voter) DeleteVoter(voterId uint) error {
	//Check if Voter exists
	redisKey := redisKeyFromId(int(voterId))
	var existingItem VoterItem
	if err := v.getItemFromRedis(redisKey, &existingItem); err != nil {
		return errors.New("voter does not exist")
	}

	if err := v.cacheClient.Del(v.context, redisKey).Err(); err != nil {
		return err
	}

	return nil
}

func (v *Voter) GetVoterPolls(voterId uint) ([]VoterHistory, error) {
	voterItem, err := v.GetVoter(voterId)
	if err != nil {
		return nil, err
	}

	return voterItem.VoteHistory, nil
}

func (v *Voter) AddVoterPoll(voterPoll VoterHistory, voterId uint) error {
	voterItem, err := v.GetVoter(voterId)
	if err != nil {
		return err
	}

	for _, vh := range voterItem.VoteHistory {
		if vh.PollId == voterPoll.PollId {
			return errors.New("poll already exists")
		}
	}

	voterItem.VoteHistory = append(voterItem.VoteHistory, voterPoll)
	if err := v.UpdateVoter(voterItem, voterId); err != nil {
		return err
	}

	return nil
}

func (v *Voter) GetVoterPoll(voterId uint, pollId uint) (VoterHistory, error) {
	voterItem, err := v.GetVoter(voterId)
	if err != nil {
		return VoterHistory{}, err
	}

	for _, vh := range voterItem.VoteHistory {
		if vh.PollId == pollId {
			return vh, nil
		}
	}

	return VoterHistory{}, errors.New("poll does not exist")
}

func (v *Voter) UpdateVoterPoll(voterPoll VoterHistory, voterId uint, pollId uint) error {
	voterItem, err := v.GetVoter(voterId)
	if err != nil {
		return err
	}

	for i, vh := range voterItem.VoteHistory {
		if vh.PollId == pollId {
			voterItem.VoteHistory[i] = voterPoll
			if err := v.UpdateVoter(voterItem, voterId); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("poll does not exist")
}

func (v *Voter) DeleteVoterPoll(voterId uint, pollId uint) error {
	voterItem, err := v.GetVoter(voterId)
	if err != nil {
		return err
	}

	for i, vh := range voterItem.VoteHistory {
		if vh.PollId == pollId {
			voterItem.VoteHistory = append(voterItem.VoteHistory[:i], voterItem.VoteHistory[i+1:]...)
			if err := v.UpdateVoter(voterItem, voterId); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("poll does not exist")
}
