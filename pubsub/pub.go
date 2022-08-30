package pubsub

type IPusher interface {
	ISubscriber
	Push() error
	Connect() error
	init(string, string)
	Reconnect() bool
}
