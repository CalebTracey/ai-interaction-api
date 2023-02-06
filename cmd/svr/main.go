package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/calebtracey/ai-interaction-api/internal"
	config "github.com/calebtracey/config-yaml"
	log "github.com/sirupsen/logrus"
)

const configPath = "config.yaml"

func main() {
	defer panicQuit()

	if svc, errs := initializeDAO(config.New(configPath)); errs != nil {
		log.Error(errs)
		panicQuit()

	} else {

		log.Fatal(listenAndServe("8080", gziphandler.GzipHandler(
			corsHandler().Handler(
				internal.Handler{
					Service: svc,
				}.InitializeRoutes(),
			)),
		))
	}
}

func panicQuit() {
	if r := recover(); r != nil {
		log.Errorf("I panicked and am quitting: %v", r)
		log.Error("I should be alerting someone...")
	}
}
