package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"jsolana/culture-360/config"
	"jsolana/culture-360/internals/app"
)

// This app is responsible to receive feedback from slack and store it providing a REST API to consume it
func main() {
	var config config.Config
	_ = envconfig.Process("app", &config)
	app, err := app.NewApplicationWithConfig(&config)
	if err != nil {
		log.WithError(err).Fatal("Unable to create the application")
		panic(err)
	}
	app.Run()
}
