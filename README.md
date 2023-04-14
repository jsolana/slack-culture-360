# Culture 360

Giving feedback is crucial for everyone's professional growth. It's important to provide feedback not only for ourselves but also for our peers. Even if the tools and methods we use change, collecting and sharing evidence of our work performance will always be the most effective way to help our team, our company, and ourselves.

If `Slack` (a communication tool used by developers and companies to share information and communicate) is your day-by-day source of truth, and to fulfill processes like 360 reviews or performance reports you spend tons of hours searching for evidences. What if you could automatize this? You could collect evidence for yourself and your peers without even thinking about it in your daily slack messages exchange.

`Culture 360` will be your friend!! ;-)

<div align="center">
  <img width="300" height="250" src="./docs/logo.webp">
</div>

Our bot will allow employees to easily create, search and collect feedback from their peers, all within `Slack`. With just a mention, the bot will store feedback related to a specific tag(s).

The best part? The bot will have an easy-to-use interface that makes the process simple for everyone involved. No more digging through old threads or struggling to find the feedback you need. With our bot, everything is in one place and easily accessible.

## How it works

We need:

- A `Slack workspace`. 
- Create the `Slack application`.
- Connecting `Culture 360` service to `Slack`.

### Slack Workspace

If you don’t already have a workspace to use, make sure to create a new one by visiting [Slack](https://slack.com/) and press `Create a new Workspace`.

Go ahead and fill all the forms, you will need to provide a name for the team or company, a new channel name, and eventually invite other teammates.

### Slack Application

The first thing we need to do is to create the `Culture 360` Slack application. Visit the [slack website](https://api.slack.com/apps/new) to create the application. Select the `From scratch` option.

You will be presented with the option to add a Name to the application (`Culture 360`) and the Workspace to allow the application to be used. You should be able to see all workspaces that you are connected to. Select the appropriate workspace.

There are many different use cases for an Application. You will be asked to select what features to add, `Culture 360` is also a bot so select the `Bot option`.

After clicking Bots you will be redirected to a Help information page, select the option to add scopes. The first thing we need to add to the application is the actual permissions to perform anything. Add the next scopes:

- `channels:history`
- `users:read` or `users.profile:read`
- `app_mentions:read`

Due `Culture360` bot listen for mentions, we require to enable the `Socket mode`. The slack events API is a way to handle events that occur in the Slack channels. There are many events, but for our bot, we want to listen to the mentions event. This means that whenever somebody mentions the bot it will receive an Event to trigger on. The events are delivered via WebSocket.

You can find all the event types available in the [documentation](https://api.slack.com/events).

The first thing you need to do is attend your Application in the [web UI](https://api.slack.com/apps/). We will activate something called `Socket Mode`, this allows the bot to connect via WebSocket. The alternative is to have the bot host a public endpoint, but then you need a domain to host on :-(.

Then we also need to add `Event Subscriptions`. You can find it in the Features tab, enter it, and activate it. Then add the `app_mentions` scope to the Event subscriptions. This will make mentions trigger a new event to the application.

The final thing we need to do is generate an Application token. Right now we only have a Bot token, but for the Events, we need an Application token.

Go into `Settings->Basic Information` and scroll down to the chapter called App-Level Tokens and press `Generate Tokens and Scope` and fill in a name for your Token. Make sure you save the Token as well by adding it to the  `.env` file as `APP_NOTIFICATION_APPTOKEN` (eg: `xapp-...`).

After adding the scopes we are ready to install the application. If you’re the owner of the workspace you can simply install it, otherwise you must request permission from an Admin.

Go into `Features->OAuth & Permissions` to obtain your Bot token and adding it to the  `.env` file as `APP_NOTIFICATION_AUTHTOKEN` (eg: `xoxb-...`).

The last step is to include the bot to the channel(s) of your interest running `/invite @culture360`.

## Connecting Culture 360 service to Slack

To build and run locally only need to check the environment variables in `.env` and run `source .env` before run `go run cmd/culture-360/main.go`.

Mainly this component requires:

- `Slack Application token`
- `Slack Auth token` (bot token)

Check `./config` to know the full available environment variables.

To test it you can use `go tool coverage`:

```batch

#Given a coverage profile produced by 'go test'
go test ./... -coverprofile=c.out

# Open a web browser displaying annotated source code
go tool cover -html=c.out

# Write out an HTML file instead of launching a web browser
go tool cover -html=c.out -o coverage.html

# Display coverage percentages to stdout for each function
go tool cover -func=c.out

```

## TODO

- Define structs: Feedback, User
- Complete the algorithm (prepare the data to be stored)
- Provided an REST API
- MySQL implementation to store / retrieve the information
- Custom metrics (principle as label)
- SOLID principles (service interfaces and implementations)
- Docker and CI/CD config
- Prepare for further steps as ML models
- Unit testing
- Improve the documentation (including instructions about how to configure app / bot). Logo and diagram
