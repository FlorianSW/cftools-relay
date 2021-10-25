package domain

type Target interface {
	Relay(e Event) error
}
