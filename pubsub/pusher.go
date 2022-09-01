package pubsub

type IPush interface {
	ISubscribe
	Push() error
	Connect() error
	init(string, string)
	Reconnect() bool
}
