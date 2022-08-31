package main

import (
	"context"
	"errors"
	"io/ioutil"
	"path/filepath"
	goRuntime "runtime"
	"strings"

	"github.com/logrusorgru/aurora" // colorable
	"go-micro.dev/v4/util/log"
	"gopkg.in/yaml.v3"

	"github.com/micro-community/stream/app"
	"github.com/micro-community/stream/runtime"
)

// Run process
func Run(ctx context.Context, configFile string) (err error) {
	_, enginePath, _, _ := goRuntime.Caller(0)
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

	go runtime.Summary.StartSummary()

	for {
		select {
		case <-ctx.Done():
			return
		}

	}
}
