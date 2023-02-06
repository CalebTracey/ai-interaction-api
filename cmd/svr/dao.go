package main

import (
	"github.com/calebtracey/ai-interaction-api/internal"
	config "github.com/calebtracey/config-yaml"
	"github.com/joho/godotenv"
)

const OpenaiApi = "openAi"

func initializeDAO() (*internal.DAO, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := config.New(configPath)

	openAiSvc, err := cfg.Service(OpenaiApi)
	if err != nil {
		return nil, err
	}

	return &internal.DAO{
		Client: openAiSvc.Client,
	}, nil
}
