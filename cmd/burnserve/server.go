package main

import (
	"fmt"
	"os"

	"github.com/mortedecai/bitburner-rpc-server/burner"
	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func main() {
	fmt.Println("Starting BitBurner Server.")
	l, _ := zap.NewDevelopment()
	logger = l.Sugar().Named("bbpusher")

	w, err := burner.NewWatcher(logger)
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
