## HATEOAS Based design

Table of contents
- [Design of APIs](#design-of-apis)
- [Golang struct changes](#golang-structure-changes)
- [New JSON Example](#new-json-structure)
- [Breif description on API usage](#brief-description-about-how-they-will-work)
- [References](#reference)


## Design of APIs
- 
- Voter API:
    - GET ``/voters``: Retrieve a list of all voters
    - POST ``/voters``: Create a new voter
    - GET ``/voters/{voterId}``: Retrieve details of a specific voter
    - GET ``/voters/{voterId}/history``: Retrieve voting history of a specific voter
    - POST ``/voters/{voterId}/vote``: Submit a vote for a specific voter
- Poll API:
    - GET ``/polls``: Retrieve a list of all polls
    - POST ``/polls``: Create a new poll
    - GET ``/polls/{pollId}``: Retrieve details of a specific poll
    - GET ``/polls/{pollId}/options``: Retrieve options available for a specific poll
- Vote API:
    - GET ``/votes``: Retrieve a list of all votes
    - POST ``/votes``: Submit a new vote
    - GET ``/votes/{voteId}``: Retrieve details of a specific vote

## Golang Structure Changes

*Voter*
```go
type Voter struct {
	VoterId string `json:"voterId"`
	Name    string `json:"name"`
	Links   struct {
		Self    Link `json:"self"`
		History Link `json:"history"`
		Polls   Link `json:"polls"`
		Vote    Link `json:"vote"`
	} `json:"_links"`
}

type VoterHistory struct {
    VoteId string `json:"voteId"`
    VoteDate time.Time `json:"voteDate"`
    Links    struct {
		Self Link `json:"self"`
        Vote Link `json:"vote"`
	} `json:"_links"`
}
```

*Poll*
```go
type Poll struct {
	PollId   string `json:"pollId"`
	Title    string `json:"title"`
	Question string `json:"question"`
	Links    struct {
		Self    Link `json:"self"`
		Options Link `json:"options"`
		Vote    Link `json:"vote"`
	} `json:"_links"`
}

type PollOptions struct {
	OptionId string `json:"optionId"`
	Text     string `json:"text"`
	Links    struct {
		Self Link `json:"self"`
	} `json:"_links"`
}
```

*Vote*
```go
type Vote struct {
	VoteId  string `json:"voteId"`
	VoterId string `json:"voterId"`
	PollId  string `json:"pollId"`
	Value   string `json:"value"`
	Links   struct {
		Self  Link `json:"self"`
		Voter Link `json:"voter"`
		Poll  Link `json:"poll"`
	} `json:"_links"`
}
```
*Link*
```go
type Link struct {
	Href string `json:"href"`
}
```

## New JSON structure

JSON will look something like this

```JSON
{
  "_links": {
    "self": { "href": "/" }
  },
  "_embedded": {
    "voters": [
      {
        "voterId": "1",
        "name": "John Doe",
        "_links": {
          "self": { "href": "/voters/1" },
          "history": { "href": "/voters/1/history" },
          "polls": { "href": "/voters/1/polls" },
          "vote": { "href": "/voters/1/vote" }
        },
        "_embedded": {
          "voteHistory": [
            {
              "voteId": "1",
              "voteDate": "2024-03-08",
              "_links": {
                "self": { "href": "/voter/1/history" },
                "vote": { "href": "/voters/1/vote" }
              }
            }
          ]
        }
      }
    ],
    "polls": [
      {
        "pollId": "1",
        "title": "Favorite Pet",
        "question": "What type of pet do you like best?",
        "_links": {
          "self": { "href": "/polls/1" },
          "options": { "href": "/polls/1/options" },
          "vote": { "href": "/polls/1/vote" }
        },
        "_embedded": {
          "options": [
            {
              "optionId": "1",
              "text": "Dogs",
              "_links": {
                "self": { "href": "/polls/1/options/1" }
              }
            },
            {
              "optionId": "2",
              "text": "Cats",
              "_links": {
                "self": { "href": "/polls/1/options/2" }
              }
            }
          ]
        }
      }
    ],
    "votes": [
      {
        "voteId": "1",
        "voterId": "1",
        "pollId": "1",
        "value": "Red",
        "_links": {
          "self": { "href": "/votes/1" },
          "voter": { "href": "/voters/1" },
          "poll": { "href": "/polls/1" }
        }
      }
    ]
  },
  "_templates": {
    "create-voter": {
      "title": "Create Voter",
      "method": "POST",
      "href": "/voters",
      "properties": {
        "name": { "type": "string", "title": "Name" }
      }
    },
    "create-poll": {
      "title": "Create Poll",
      "method": "POST",
      "href": "/polls",
      "properties": {
        "title": { "type": "string", "title": "Title" },
        "question": { "type": "string", "title": "Question" }
      }
    },
    "create-poll-options": {
      "title": "Create Poll Options",
      "method": "POST",
      "href": "/polls/{id}/options",
      "properties": {
        "text": { "type": "string", "title": "Text" }
      }
    },
    "create-vote": {
      "title": "Create Vote",
      "method": "POST",
      "href": "/votes",
      "properties": {
        "voterId": { "type": "string", "title": "Voter ID" },
        "pollId": { "type": "string", "title": "Poll ID" },
        "value": { "type": "string", "title": "Value" }
      }
    }
  }
}

```
We are using HAL specification for the JSON structure, where
- `_links` at the top defines the context
- `_embedded` defines resources which contains entities
- `entities` are vote, poll and voter
- `_templates` defines the templates that can be used to consue the API for example `create-voter` template can be used to create a new voter by using HAL SDKs

## Brief Description about how they will work

The flow of API is still same i.e. you need a voter and poll to be created first before a vote is cast. However with hypermedia in the responses, now API consumer have API media and template as a guide to get started with API. Lets go through one sample flow

1. Create Voter which will return the links to self, history, polls and vote.
2. Using polls new poll is created to ask voters about their favourite pet. That will return the links to self, options and vote. And using options link pet options are added
3. Now using vote API, vote can be cast.

This is how high level flow looks like and API response itself provide all the information for next APIs to call.

## Reference

- https://en.wikipedia.org/wiki/HATEOAS
- https://www.mscharhag.com/api-design/hypermedia-rest
- https://restfulapi.net/hateoas/
- https://en.wikipedia.org/wiki/Hypertext_Application_Language
- https://stateless.co/hal_specification.html
- https://github.com/byteclubfr/js-hal
