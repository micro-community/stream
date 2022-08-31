package pubsub

type Subscribe interface {
	IChannel
	receive(string, IChannel) error
	IsPlaying() bool
	PlayRaw()
	PlayBlock(byte)
	PlayFLV()
	Stop()
}
