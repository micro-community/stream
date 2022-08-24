package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"go-micro.dev/v4/util/log"
	"gopkg.in/yaml.v3"

	// colorable
	"github.com/logrusorgru/aurora"
	"github.com/micro-community/stream/app"
	"github.com/micro-community/stream/engine"
)

// Run process
func Run(ctx context.Context, configFile string) (err error) {

	_, enginePath, _, _ := runtime.Caller(0)
	if parts := strings.Split(filepath.Dir(enginePath), "@"); len(parts) > 1 {
		app.Version = parts[len(parts)-1]
	}
	rawConfigData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Info(aurora.Red("read config file error:"), err)
		return errors.New("read config file error")
	}

	//get config
	var conf app.Configuration
	if err := yaml.Unmarshal(rawConfigData, &conf); err != nil {
		log.Error("parsing yml files error:,system will use default", err)
	}

	log.Info(aurora.Green("start stream server"), aurora.BrightBlue(app.Version))

	go engine.Summary.StartSummary()

	var pluginConfigs map[string]interface{}
	for name, config := range engine.Plugins {
		if cfg, ok := pluginConfigs[name]; ok {
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
	return
}
