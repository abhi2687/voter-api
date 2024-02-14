package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/abhi2687/voter-api/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	hostFlag           string
	portFlag           uint
	app                *fiber.App
	voterHandler       *api.VoterAPI
	err                error
	startTime          = time.Now()
	successfulRequests = 0
	failedRequests     = 0
)

func main() {
	processCommandLineFlag()
	initializeAppUsingFiber()
	initializeVoterAPIHandler()
	registerHandlers()
	StartServer()
}

// Middleware function to count successful requests
func countSuccessfulRequests(c *fiber.Ctx) error {
	err := c.Next()
	if err == nil {
		successfulRequests++
	}
	return err
}

// Middleware function to count failed requests
func countFailedRequests(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		failedRequests++
	}
	return err
}

func initializeVoterAPIHandler() {
	voterHandler, err = api.New()
	if err != nil {
		fmt.Printf("Error creating voter handler: %v\n", err)
		os.Exit(1)
	}
}

func registerHandlers() {
	app.Get("/voters/health", HealthCheck)
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

func initializeAppUsingFiber() {
	app = fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(countSuccessfulRequests, countFailedRequests)
}

func processCommandLineFlag() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")
	flag.Parse()
}

func StartServer() {
	// Start the server
	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	log.Println("Starting server on ", serverPath)
	app.Listen(serverPath)
}

func HealthCheck(c *fiber.Ctx) error {
	uptime := time.Since(startTime).Seconds()
	return c.Status(http.StatusOK).
		JSON(fiber.Map{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             uptime,
			"successfulRequests": successfulRequests,
			"failedRequest":      failedRequests,
		})
}
