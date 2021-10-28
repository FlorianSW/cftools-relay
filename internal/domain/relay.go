package domain

type Target interface {
	Relay(e Event, f *Filter) error
}
