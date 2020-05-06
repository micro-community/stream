package record

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/micro-community/x-streaming/engine"
	"github.com/micro-community/x-streaming/engine/avformat"
)

//FlvFile for flv container
type FlvFile struct {
	engine.Publisher
}

//PublishFlvFile pub flv file streaming
func PublishFlvFile(streamPath string) error {
	if file, err := os.Open(filepath.Join(config.Path, streamPath+".flv")); err == nil {
		stream := FlvFile{}
		if stream.Publish(streamPath) {
			stream.Type = "FlvFile"
			defer stream.Close()
			stream.UseTimestamp = true
			file.Seek(int64(len(avformat.FLVHeader)), io.SeekStart)
			var lastTime uint32
			for {
				if t, timestamp, payload, err := avformat.ReadFLVTag(file); err == nil {
					switch t {
					case avformat.FLV_TAG_TYPE_AUDIO:
						stream.PushAudio(timestamp, payload)
					case avformat.FLV_TAG_TYPE_VIDEO:
						if timestamp != 0 {
							if lastTime == 0 {
								lastTime = timestamp
							}
						}
						stream.PushVideo(timestamp, payload)
						time.Sleep(time.Duration(timestamp-lastTime) * time.Millisecond)
						lastTime = timestamp
					}
				} else {
					return err
				}
			}
		} else {
			return errors.New("Bad Name")
		}
	} else {
		return err
	}
}
