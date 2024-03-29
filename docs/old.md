# Old implementation

This is the seed, processing mentions, commands using slack-go and socket mode:

```golang

package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "os"
    "strings"
    "time"

    "github.com/joho/godotenv"
    "github.com/slack-go/slack"
    "github.com/slack-go/slack/slackevents"
    "github.com/slack-go/slack/socketmode"
)

func main() {

    // Load Env variables from .dot file
    godotenv.Load(".env")

    token := os.Getenv("SLACK_AUTH_TOKEN")
    appToken := os.Getenv("SLACK_APP_TOKEN")
    // Create a new client to slack by giving token
    // Set debug to true while developing
    // Also add a ApplicationToken option to the client
    client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
    // go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
    socketClient := socketmode.New(
        client,
        socketmode.OptionDebug(true),
        // Option to set a custom logger
        socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
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
                log.Println("Shutting down socketmode listener")
                return
            case event := <-socketClient.Events:
                // We have a new Events, let's type switch the event
                // Add more use cases here if you want to listen to other events.
                switch event.Type {
                case socketmode.EventTypeConnecting:
                    fmt.Println("Connecting to Slack with Socket Mode...")
                case socketmode.EventTypeConnectionError:
                    fmt.Println("Connection failed. Retrying later...")
                case socketmode.EventTypeConnected:
                    fmt.Println("Connected to Slack with Socket Mode.")
                case socketmode.EventTypeHello:
                    fmt.Println("Hello received...")
                // handle EventAPI events
                case socketmode.EventTypeEventsAPI:
                    // The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
                    eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)

                    if !ok {
                        log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
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
                    // Handle Slash Commands
                case socketmode.EventTypeSlashCommand:
                    // Just like before, type cast to the correct event type, this time a SlashEvent
                    command, ok := event.Data.(slack.SlashCommand)
                    if !ok {
                        log.Printf("Could not type cast the message to a SlashCommand: %v\n", command)
                        continue
                    }
                    // Dont forget to acknowledge the request
                    socketClient.Ack(*event.Request)
                    // handleSlashCommand will take care of the command
                    err := handleSlashCommand(command, client)
                    if err != nil {
                        log.Fatal(err)
                    }
                default:
                    fmt.Printf("Event received: %v\n", event)
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
        return errors.New("unsupported event type")
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
    // Check if the user said Hello to the bot
    text := strings.ToLower(event.Text)

    // Create the attachment and assigned based on the message
    attachment := slack.Attachment{}
    // Add Some default context like user who mentioned the bot
    attachment.Fields = []slack.AttachmentField{
        {
            Title: "Date",
            Value: time.Now().String(),
        }, {
            Title: "Initializer",
            Value: user.Name,
        },
    }
    if strings.Contains(text, "hello") {
        // Greet the user
        attachment.Text = fmt.Sprintf("Hello %s", user.Name)
        attachment.Pretext = "Greetings"
        attachment.Color = "#4af030"
    } else {
        // Send a message to the user
        attachment.Text = fmt.Sprintf("How can I help you %s?", user.Name)
        attachment.Pretext = "How can I be of service"
        attachment.Color = "#3d3d3d"
    }
    // Send the message to the channel
    // The Channel is available in the event message
    _, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
    if err != nil {
        return fmt.Errorf("failed to post message: %w", err)
    }
    return nil
    }

    // handleSlashCommand will take a slash command and route to the appropriate function
    func handleSlashCommand(command slack.SlashCommand, client *slack.Client) error {
        // We need to switch depending on the command
        switch command.Command {
        case "/fb":
            // This was a feedback command, so pass it along to the proper function
            return handleFeedbackCommand(command, client)
        }

        return nil
    }

    // handleFeedbackCommand will take care of /hello submissions
    func handleFeedbackCommand(command slack.SlashCommand, client *slack.Client) error {
        fmt.Printf("%s command %v\n", command.UserName, command)
        _, _, err := client.PostMessage(command.ChannelID, slack.MsgOptionText("Ok!", false))
        if err != nil {
            return fmt.Errorf("failed to post message: %w", err)
        }
        return nil
    }

```
