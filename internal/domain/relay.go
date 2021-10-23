package domain

type Target interface {
	Relay(e WebhookEvent) error
}
