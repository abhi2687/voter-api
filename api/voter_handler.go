package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/abhi2687/voter-api/db"
	"github.com/gofiber/fiber/v2"
)

type VoterAPI struct {
	db *db.VoterList
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.New()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{db: dbHandler}, nil
}

func (v *VoterAPI) AddVoter(c *fiber.Ctx) error {
	var voter db.Voter
	fmt.Println("Request body: ", string(c.Body()))
	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error parsing request body: ", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := v.db.AddVoter(voter)
	if err != nil {
		log.Println("Error adding voter: ", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(voter)
}

func (v *VoterAPI) GetVoter(c *fiber.Ctx) error {
	voterIdStr := c.Params("id")
	voterId, err := strconv.ParseUint(voterIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing voterId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	voter, err := v.db.GetVoter(uint(voterId))
	if err != nil {
		log.Println("Error getting voter: ", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(voter)
}

func (v *VoterAPI) GetAllVoters(c *fiber.Ctx) error {
	voters := v.db.GetAllVoters()
	return c.Status(http.StatusOK).JSON(voters)
}

func (v *VoterAPI) DeleteAllVoters(c *fiber.Ctx) error {
	v.db.DeleteAllVoters()
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (v *VoterAPI) UpdateVoter(c *fiber.Ctx) error {
	var voter db.Voter
	voterIdStr := c.Params("id")
	voterId, err := strconv.ParseUint(voterIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing voterId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error parsing request body", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = v.db.UpdateVoter(voter, uint(voterId))
	if err != nil {
		log.Println("Error updating voter: ", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (v *VoterAPI) DeleteVoter(c *fiber.Ctx) error {
	voterIdStr := c.Params("id")
	voterId, err := strconv.ParseUint(voterIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing voterId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = v.db.DeleteVoter(uint(voterId))
	if err != nil {
		log.Println("Error deleting voter: ", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (v *VoterAPI) GetVoterPolls(c *fiber.Ctx) error {
	voterIdStr := c.Params("id")
	voterId, err := strconv.ParseUint(voterIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing voterId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	voterPolls, err := v.db.GetVoterPolls(uint(voterId))
	if err != nil {
		log.Println("Error getting voter: ", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(voterPolls)
}

func (v *VoterAPI) AddVoterPoll(c *fiber.Ctx) error {
	var voterPoll db.VoterHistory
	voterIdStr := c.Params("id")
	voterId, err := strconv.ParseUint(voterIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing voterId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.BodyParser(&voterPoll); err != nil {
		log.Println("Error parsing request body: ", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = v.db.AddVoterPoll(voterPoll, uint(voterId))
	if err != nil {
		log.Println("Error adding voter poll: ", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"status": "ok"})
}

func (v *VoterAPI) GetVoterPoll(c *fiber.Ctx) error {
	voterIdStr := c.Params("id")
	voterId, err := strconv.ParseUint(voterIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing voterId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	pollIdStr := c.Params("pollid")
	pollId, err := strconv.ParseUint(pollIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing pollId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	voterPoll, err := v.db.GetVoterPoll(uint(voterId), uint(pollId))
	if err != nil {
		log.Println("Error getting voter poll: ", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(voterPoll)
}

func (v *VoterAPI) UpdateVoterPoll(c *fiber.Ctx) error {
	var voterPoll db.VoterHistory
	voterIdStr := c.Params("id")
	voterId, err := strconv.ParseUint(voterIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing voterId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	pollIdStr := c.Params("pollid")
	pollId, err := strconv.ParseUint(pollIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing pollId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.BodyParser(&voterPoll); err != nil {
		log.Println("Error parsing request body: ", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = v.db.UpdateVoterPoll(voterPoll, uint(voterId), uint(pollId))
	if err != nil {
		log.Println("Error updating voter poll: ", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (v *VoterAPI) DeleteVoterPoll(c *fiber.Ctx) error {
	voterIdStr := c.Params("id")
	voterId, err := strconv.ParseUint(voterIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing voterId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	pollIdStr := c.Params("pollid")
	pollId, err := strconv.ParseUint(pollIdStr, 10, 32)
	if err != nil {
		log.Println("Error parsing pollId", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = v.db.DeleteVoterPoll(uint(voterId), uint(pollId))
	if err != nil {
		log.Println("Error deleting voter poll: ", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
}
