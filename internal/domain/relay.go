package domain

type Target interface {
	Relay(message string, e WebhookEvent) error
}
