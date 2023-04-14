package app

import (
	"context"
	"errors"
	"fmt"
	"jsolana/culture-360/config"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Application struct {
	config *config.Config
}

// Returns a new Application instance from config
func NewApplicationWithConfig(cfg *config.Config) (*Application, error) {
	if cfg == nil {
		return nil, errors.New("configuration can't be nil")
	}

	return NewApplication(cfg)
}

// Returns a new Application instance
func NewApplication(cfg *config.Config) (*Application, error) {
	if cfg == nil {
		return nil, errors.New("configuration can't be nil")
	}

	err := configureLogging(cfg)

	if err != nil {
		return nil, err
	}

	configureMetrics(cfg) // ignore errors

	defer log.Infof("%s initialized", cfg.Name)
	return &Application{
		config: cfg,
	}, nil
}

// Run the application
// 1. Read the pending invitations
// 2. Transform to render reports
// 3. Notify the rendered reports
// Register metrics
func (app *Application) Run() {
	// Create a new client to slack by giving token
	// Set debug to true while developing
	// Also add a ApplicationToken option to the client
	client := slack.New(app.config.Notification.AuthToken, slack.OptionDebug(false), slack.OptionAppLevelToken(app.config.Notification.AppToken))
	// go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(false),
	)

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		// Create a for loop that selects either the context cancellation or the events incomming
		for {
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				log.Info("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				// We have a new Events, let's type switch the event
				// Add more use cases here if you want to listen to other events.
				switch event.Type {
				// handle EventAPI events
				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)

					if !ok {
						log.Errorf("Could not type cast the event to the EventsAPIEvent: %v", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socketClient.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					err := handleEventMessage(eventsAPIEvent, client)
					if err != nil {
						// Replace with actual err handeling
						log.Fatal(err)
					}

				default:
					log.Debugf("Event received: %v", event)
				}

			}

		}
	}(ctx, client, socketClient)

	socketClient.Run()
}

// handleEventMessage will take an event and handle it properly based on the type of event
func handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := handleAppMentionEvent(ev, client)
			if err != nil {
				return err
			}
		}
	default:
		log.Infof("Event type received: %v", event)
		return nil
	}
	return nil
}

// handleAppMentionEvent is used to take care of the AppMentionEvent when the bot is mentioned
func handleAppMentionEvent(event *slackevents.AppMentionEvent, client *slack.Client) error {

	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	// Check the text
	text := strings.ToLower(event.Text)
	//TODO extract user, channel, create the link, date...
	log.Infof("Mention >>>>>> %v received from user %s", text, user.Name)
	return nil
}

// Configure  logs
func configureLogging(cfg *config.Config) error {
	log.SetFormatter(cfg.NewLoggingFormatter())
	level, err := cfg.LoggingLevel()
	if err != nil {
		return err
	}

	log.SetLevel(level)
	return nil
}

// configure metrics (register metrics and run a http listener for prometheus)
func configureMetrics(cfg *config.Config) {
	if cfg.Metrics.Prometheus.Enabled {
		go func() {
			prometheus.MustRegister(BatchHistogram)
			http.Handle(cfg.Metrics.Prometheus.HTTPPath, promhttp.Handler())
			log.Infof("Enabling prometheus %s endpoint and port %s", cfg.Metrics.Prometheus.HTTPPath, cfg.Metrics.Prometheus.Port)
			if err := http.ListenAndServe(fmt.Sprintf(":"+cfg.Metrics.Prometheus.Port), nil); err != nil && err != http.ErrServerClosed {
				log.Errorf("An error occurs initializing Prometheus: %v", err)
				//	ignoring
			}
		}()
	}
}
