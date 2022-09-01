package pubsub

type PublishOption struct {
	PubAudio          bool
	PubVideo          bool
	KickExist         bool // 是否踢掉已经存在的发布者
	PublishTimeout    int  // 发布无数据超时
	WaitCloseTimeout  int  // 延迟自动关闭（等待重连）
	DelayCloseTimeout int  // 延迟自动关闭（无订阅时）
}

type SubscribeOption struct {
	SubAudio    bool
	SubVideo    bool
	LiveMode    bool // 实时模式：追赶发布者进度，在播放首屏后等待发布者的下一个关键帧，然后调到该帧。
	IFrameOnly  bool // 只要关键帧
	WaitTimeout int  // 等待流超时
}

// PushOption for push stream
type PushOption struct {
	RePush   int               // 断开后自动重推,0 表示不自动重推，-1 表示无限重推，高于0 的数代表最大重推次数
	PushList map[string]string // 自动推流列表
}

// PullOption for pull stream
type PullOption struct {
	RePull          int               // 断开后自动重拉,0 表示不自动重拉，-1 表示无限重拉，高于0 的数代表最大重拉次数
	PullOnStart     bool              // 启动时拉流
	PullOnSubscribe bool              // 订阅时自动拉流
	PullList        map[string]string // 自动拉流列表，以streamPath为key，url为value
}
