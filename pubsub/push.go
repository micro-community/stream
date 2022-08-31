package pubsub

type IPush interface {
	Subscribe
	Push() error
	Connect() error
	init(string, string)
	Reconnect() bool
}
