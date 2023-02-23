package main

import (
	"fmt"
	"os"

	"github.com/mortedecai/bitburner-rpc-server/burner"
	"github.com/mortedecai/go-go-gadgets/env"
	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func main() {
	fmt.Println("Starting BitBurner Server.")
	l, _ := zap.NewDevelopment()
	logger = l.Sugar().Named("bbpusher")

	var token string
	if t, found := env.GetWithDefault("BB_API_TOKEN", ""); !found {
		logger.Errorw("main", "API Token", "NOT FOUND", "EnvVar", "BB_API_TOKEN")
		return
	} else {
		token = t
		logger.Infow("main", "Token", token)
	}

	w, err := burner.NewWatcher(logger, token)
	if err != nil {
		logger.Fatalw("main", "Create Watcher Error", err)
	}
	w.Start()

	if len(os.Args) == 1 {
		if cwd, err := os.Getwd(); err == nil {
			w.AddDirectory(cwd)
		} else {
			logger.Fatalw("main", "Args 0 Watch Error", err)
		}
	} else {
		const argString = "os.Args[%d]"
		for i, v := range os.Args[1:] {
			if finfo, err := os.Stat(v); err == nil {
				if finfo.IsDir() {
					logger.Infow("main", fmt.Sprintf(argString, i), v, "Directory?", "true")
					w.AddDirectory(v)
				} else {
					logger.Infow("main", fmt.Sprintf(argString, i), v, "Directory?", "false")
				}
			}
		}
	}
	<-make(chan struct{})
}
