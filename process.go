package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"go-micro.dev/v4/util/log"

	// colorable
	"github.com/BurntSushi/toml"
	"github.com/logrusorgru/aurora"
	"github.com/micro-community/stream/engine"
)

// Run engine
func Run(ctx context.Context, configFile string) (err error) {

	_, enginePath, _, _ := runtime.Caller(0)
	if parts := strings.Split(filepath.Dir(enginePath), "@"); len(parts) > 1 {
		engine.Version = parts[len(parts)-1]
	}
	if engine.ConfigRaw, err = ioutil.ReadFile(configFile); err != nil {
		engine.Print(aurora.Red("read config file error:"), err)
		return
	}
	engine.Print(aurora.Green("start monibuca"), aurora.BrightBlue(engine.Version))
	go engine.Summary.StartSummary()
	var cg map[string]interface{}
	if _, err = toml.Decode(string(engine.ConfigRaw), &cg); err == nil {
		if cfg, ok := cg["Monibuca"]; ok {
			b, _ := json.Marshal(cfg)
			if err = json.Unmarshal(b, engine.Config); err != nil {
				log.Error(err)
			}
		}
		for name, config := range engine.Plugins {
			if cfg, ok := cg[name]; ok {
				b, _ := json.Marshal(cfg)
				if err = json.Unmarshal(b, config.Config); err != nil {
					log.Error(err)
					continue
				}
			} else if config.Config != nil {
				continue
			}
			if config.Run != nil {
				go config.Run()
			}
		}
	} else {
		engine.Print(aurora.Red("decode config file error:"), err)
	}
	return
}
