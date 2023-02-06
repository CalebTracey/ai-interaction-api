package main

import (
	"github.com/calebtracey/ai-interaction-api/internal"
	config "github.com/calebtracey/config-yaml"
	"github.com/joho/godotenv"
)

const OpenaiApi = "openAi"

func initializeDAO(appConfig *config.Config) (*internal.DAO, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	openAiSvc, err := appConfig.Service(OpenaiApi)
	if err != nil {
		return nil, err
	}

	return &internal.DAO{
		Client: openAiSvc.Client,
	}, nil
}
