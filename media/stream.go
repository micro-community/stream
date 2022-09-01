package media

type IStream interface {
	AddTrack(Track)
	RemoveTrack(Track)
	IsClosed() bool
	SSRC() uint32
	Receive(any) bool
}
