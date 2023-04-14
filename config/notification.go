package config

// SlackConfiguration is the configuration of the notifications using slack
type SlackConfiguration struct {
	// Token provided by slack for authentication
	// Environment variable: APP_NOTIFICATION_AUTHTOKEN
	AuthToken string
	// Token provided by slack for app authentication
	// Environment variable: APP_NOTIFICATION_APPTOKEN
	AppToken string
}
