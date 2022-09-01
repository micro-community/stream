package pubsub

// Publish of Publisher
type IPublish interface {
	IChannel
	//get option
	GetOption() PublishOption
	//receive stream from channel
	receive(string, IChannel, PublishOption) error
	//Get channel for publishing
	GetChannel() *Channel[PublishOption]

	getAudioTrack() media.AudioTrack
	getVideoTrack() media.VideoTrack
}
