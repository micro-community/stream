package engine

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/micro-community/stream/app"
	"github.com/micro-community/stream/codecs"
)

var (
	streamCollection = Collection{}
)

// Collection 对sync.Map的包装
type Collection struct {
	sync.Map
}

// FindStream 根据流路径查找流
func FindStream(streamPath string) *Stream {
	if s, ok := streamCollection.Load(streamPath); ok {
		return s.(*Stream)
	}
	return nil
}

// GetStream 根据流路径获取流，如果不存在则创建一个新的
func GetStream(streamPath string) (result *Stream) {
	item, loaded := streamCollection.LoadOrStore(streamPath, &Stream{
		Subscribers: make(map[string]*Subscriber),
		Control:     make(chan interface{}),
		AVRing:      NewRing(app.Config.RingSize),
		StreamInfo: StreamInfo{
			StreamPath:     streamPath,
			SubscriberInfo: make([]*SubscriberInfo, 0),
			HasVideo:       true,
			HasAudio:       true,
			EnableAudio:    &app.Config.EnableAudio,
			EnableVideo:    &app.Config.EnableVideo,
		},
		WaitPub: make(chan struct{}),
	})
	result = item.(*Stream)
	if !loaded {
		Summary.Streams = append(Summary.Streams, &result.StreamInfo)
		result.Context, result.Cancel = context.WithCancel(context.Background())
		if app.Config.EnableVideo {
			result.EnableVideo = &result.HasVideo
		}
		if app.Config.EnableAudio {
			result.EnableAudio = &result.HasAudio
		}
		go result.Run()
	}
	return
}

// Stream 流定义
type Stream struct {
	context.Context
	*Publisher
	StreamInfo   //可序列化，供后台查看的数据
	Control      chan interface{}
	Cancel       context.CancelFunc
	Subscribers  map[string]*Subscriber // 订阅者
	VideoTag     *codecs.AVPacket       // 每个视频包都是这样的结构,区别在于Payload的大小.FMS在发送AVC sequence header,需要加上 VideoTags,这个tag 1个字节(8bits)的数据
	AudioTag     *codecs.AVPacket       // 每个音频包都是这样的结构,区别在于Payload的大小.FMS在发送AAC sequence header,需要加上 AudioTags,这个tag 1个字节(8bits)的数据
	FirstScreen  *Ring                  //最近的关键帧位置，首屏渲染
	AVRing       *Ring                  //数据环
	WaitPub      chan struct{}          //用于订阅和等待发布者
	UseTimestamp bool                   //是否采用数据包中的时间戳
	SPS          []byte
	PPS          []byte
}

// StreamInfo 流可序列化信息，用于控制台显示
type StreamInfo struct {
	StreamPath     string
	StartTime      time.Time
	SubscriberInfo []*SubscriberInfo
	Type           string
	VideoInfo      struct {
		PacketCount int
		CodecID     byte
		SPSInfo     codecs.SPSInfo
		BPS         int
		lastIndex   int
		GOP         int //关键帧间隔
	}
	AudioInfo struct {
		PacketCount int
		SoundFormat byte //4bit
		SoundRate   int  //2bit
		SoundSize   byte //1bit
		SoundType   byte //1bit
		lastIndex   int
		BPS         int
	}
	HasAudio    bool
	HasVideo    bool
	EnableVideo *bool
	EnableAudio *bool
}

// UnSubscribeCmd 取消订阅命令
type UnSubscribeCmd struct {
	*Subscriber
}

// SubscribeCmd 订阅流命令
type SubscribeCmd struct {
	*Subscriber
}

// ChangeStreamCmd 切换流命令
type ChangeStreamCmd struct {
	*Subscriber
	NewStream *Stream
}

func (r *Stream) onClosed() {
	Print(aurora.Yellow("Stream destroyed :"), aurora.BrightCyan(r.StreamPath))
	streamCollection.Delete(r.StreamPath)
	for i, val := range Summary.Streams {
		if val == &r.StreamInfo {
			Summary.Streams = append(Summary.Streams[:i], Summary.Streams[i+1:]...)
			break
		}
	}
	OnStreamClosedHooks.Trigger(r)
}

// Subscribe 订阅流
func (r *Stream) Subscribe(s *Subscriber) {
	s.Stream = r
	if r.Err() == nil {
		s.SubscribeTime = time.Now()
		Print(aurora.Sprintf(aurora.Yellow("subscribe :%s %s,to Stream %s"), aurora.Blue(s.Type), aurora.Cyan(s.ID), aurora.BrightCyan(r.StreamPath)))
		s.Context, s.Cancel = context.WithCancel(r)
		s.Control <- &SubscribeCmd{s}
	}
}

// UnSubscribe 取消订阅流
func (r *Stream) UnSubscribe(s *Subscriber) {
	if r.Err() == nil {
		r.Control <- &UnSubscribeCmd{s}
	}
}

// Run 流运行
func (r *Stream) Run() {
	Print(aurora.Green("Stream create:"), aurora.BrightCyan(r.StreamPath))
	defer r.onClosed()
	for {
		select {
		case <-r.Done():
			return
		case s := <-r.Control:
			switch v := s.(type) {
			case *UnSubscribeCmd:
				if _, ok := r.Subscribers[v.ID]; ok {
					delete(r.Subscribers, v.ID)
					for i, val := range r.SubscriberInfo {
						if val == &v.SubscriberInfo {
							r.SubscriberInfo = append(r.SubscriberInfo[:i], r.SubscriberInfo[i+1:]...)
							break
						}
					}
					OnUnSubscribeHooks.Trigger(v.Subscriber)
					Print(aurora.Sprintf(aurora.Yellow("%s subscriber %s removed remains:%d"), aurora.BrightCyan(r.StreamPath), aurora.Cyan(v.ID), aurora.Blue(len(r.SubscriberInfo))))
					if len(r.SubscriberInfo) == 0 && r.Publisher == nil {
						r.Cancel()
					}
				}
			case *SubscribeCmd:
				//防止重复添加
				if _, ok := r.Subscribers[v.ID]; !ok {
					r.Subscribers[v.ID] = v.Subscriber
					r.SubscriberInfo = append(r.SubscriberInfo, &v.SubscriberInfo)
					Print(aurora.Sprintf(aurora.Yellow("%s subscriber %s added remains:%d"), aurora.BrightCyan(r.StreamPath), aurora.Cyan(v.ID), aurora.Blue(len(r.SubscriberInfo))))
					OnSubscribeHooks.Trigger(v.Subscriber)
				}
			case *ChangeStreamCmd:
				if _, ok := v.NewStream.Subscribers[v.ID]; !ok {
					delete(r.Subscribers, v.ID)
					v.NewStream.Subscribe(v.Subscriber)
					if len(r.SubscriberInfo) == 0 && r.Publisher == nil {
						r.Cancel()
					}
				}
			}
		}
	}
}

// PushAudio 来自发布者推送的音频
func (r *Stream) PushAudio(timestamp uint32, payload []byte) {
	payloadLen := len(payload)
	audio := r.AVRing
	audio.Type = codecs.FLV_TAG_TYPE_AUDIO
	audio.Timestamp = timestamp
	audio.Payload = payload
	audio.IsKeyFrame = false
	audio.IsSequence = false

	if payloadLen < 4 {
		return
	}
	if payload[0] == 0xFF && (payload[1]&0xF0) == 0xF0 {
		//将ADTS转换成ASC
		r.AudioInfo.SoundFormat = 10
		r.AudioInfo.SoundRate = codecs.SamplingFrequencies[(payload[2]&0x3c)>>2]
		r.AudioInfo.SoundType = ((payload[2] & 0x1) << 2) | ((payload[3] & 0xc0) >> 6)
		r.AudioTag = audio.ADTS2ASC()
	} else if r.AudioTag == nil {
		audio.IsSequence = true
		if payloadLen < 5 {
			return
		}
		r.AudioTag = audio.AVPacket.Clone()
		tmp := payload[0]                                                      // 第一个字节保存着音频的相关信息
		if r.AudioInfo.SoundFormat = tmp >> 4; r.AudioInfo.SoundFormat == 10 { //真的是AAC的话，后面有一个字节的详细信息
			//0 = AAC sequence header，1 = AAC raw。
			if aacPacketType := payload[1]; aacPacketType == 0 {
				config1 := payload[2]
				config2 := payload[3]
				//audioObjectType = (config1 & 0xF8) >> 3
				// 1 AAC MAIN 	ISO/IEC 14496-3 subpart 4
				// 2 AAC LC 	ISO/IEC 14496-3 subpart 4
				// 3 AAC SSR 	ISO/IEC 14496-3 subpart 4
				// 4 AAC LTP 	ISO/IEC 14496-3 subpart 4
				r.AudioInfo.SoundRate = codecs.SamplingFrequencies[((config1&0x7)<<1)|(config2>>7)]
				r.AudioInfo.SoundType = (config2 >> 3) & 0x0F //声道
				//frameLengthFlag = (config2 >> 2) & 0x01
				//dependsOnCoreCoder = (config2 >> 1) & 0x01
				//extensionFlag = config2 & 0x01
			}
		} else {
			r.AudioInfo.SoundRate = codecs.SoundRate[(tmp&0x0c)>>2] // 采样率 0 = 5.5 kHz or 1 = 11 kHz or 2 = 22 kHz or 3 = 44 kHz
			r.AudioInfo.SoundSize = (tmp & 0x02) >> 1               // 采样精度 0 = 8-bit samples or 1 = 16-bit samples
			r.AudioInfo.SoundType = tmp & 0x01                      // 0 单声道，1立体声
		}
		return
	}
	if !r.UseTimestamp {
		audio.Timestamp = uint32(time.Since(r.StartTime) / time.Millisecond)
	}
	lastTimestamp := audio.GetAt(r.AudioInfo.lastIndex).Timestamp
	if lastTimestamp > 0 && lastTimestamp != audio.Timestamp {
		r.AudioInfo.BPS = payloadLen * 1000 / int(audio.Timestamp-lastTimestamp)
	}
	r.AudioInfo.PacketCount++
	audio.Number = r.AudioInfo.PacketCount
	r.AudioInfo.lastIndex = audio.Index
	audio.NextW()
}
func (r *Stream) setH264Info(video *Ring) {
	r.VideoTag = video.AVPacket.Clone()
	if r.VideoInfo.CodecID != 7 {
		return
	}
	info := codecs.AVCDecoderConfigurationRecord{}
	//0:codec,1:IsAVCSequence,2~4:compositionTime
	if _, err := info.Unmarshal(video.Payload[5:]); err == nil {
		r.VideoInfo.SPSInfo, err = codecs.ParseSPS(info.SequenceParameterSetNALUnit)
		if err != nil {
			return
		}
	}
}

func (r *Stream) WriteSPS(sps []byte) {
	lenSPS := len(sps)
	r.SPS = sps
	if r.VideoTag == nil {
		r.VideoTag = codecs.NewAVPacket(codecs.FLV_TAG_TYPE_VIDEO)
		r.VideoTag.IsSequence = true
		r.VideoTag.IsKeyFrame = true
		r.VideoTag.Payload = append(r.VideoTag.Payload, codecs.RTMP_AVC_HEAD...)
	}
	r.VideoInfo.SPSInfo, _ = codecs.ParseSPS(sps)
	copy(r.VideoTag.Payload[6:], sps[1:4])
	r.VideoTag.Payload = append(append(r.VideoTag.Payload[:10], 0xE1, byte(lenSPS>>8), byte(lenSPS)), sps...)
}
func (r *Stream) WritePPS(pps []byte) {
	lenPPS := len(pps)
	r.PPS = pps
	r.VideoTag.Payload = append(append(r.VideoTag.Payload, 0x01, byte(lenPPS>>8), byte(lenPPS)), pps...)
}

// PushVideo 来自发布者推送的视频
func (r *Stream) PushVideo(timestamp uint32, payload []byte) {
	payloadLen := len(payload)
	if payloadLen < 3 {
		return
	}
	video := r.AVRing
	video.Type = codecs.FLV_TAG_TYPE_VIDEO
	video.Timestamp = timestamp
	video.Payload = payload
	videoFrameType := payload[0] >> 4       // 帧类型 4Bit, H264一般为1或者2
	r.VideoInfo.CodecID = payload[0] & 0x0f // 编码类型ID 4Bit, JPEG, H263, AVC...
	video.IsSequence = videoFrameType == 1 && payload[1] == 0
	video.IsKeyFrame = videoFrameType == 1 || videoFrameType == 4
	r.VideoInfo.PacketCount++
	video.Number = r.VideoInfo.PacketCount
	if r.VideoTag == nil {
		if video.IsSequence {
			r.setH264Info(video)
		} else {
			log.Println("no AVCSequence")
		}
	} else {
		//更换AVCSequence
		if video.IsSequence {
			r.setH264Info(video)
		}
		if video.IsKeyFrame {
			if r.FirstScreen == nil {
				defer close(r.WaitPub)
				r.FirstScreen = video.Clone()
			} else {
				oldNumber := r.FirstScreen.Number
				r.FirstScreen.GoTo(video.Index)
				r.VideoInfo.GOP = r.FirstScreen.Number - oldNumber
			}
		}
		if !r.UseTimestamp {
			video.Timestamp = uint32(time.Since(r.StartTime) / time.Millisecond)
		}
		lastTimestamp := video.GetAt(r.VideoInfo.lastIndex).Timestamp
		if lastTimestamp > 0 && lastTimestamp != video.Timestamp {
			r.VideoInfo.BPS = payloadLen * 1000 / int(video.Timestamp-lastTimestamp)
		}
		r.VideoInfo.lastIndex = video.Index
		video.NextW()
	}
}
