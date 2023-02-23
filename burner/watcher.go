package burner

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type Watcher struct {
	spy     *fsnotify.Watcher
	logger  *zap.SugaredLogger
	handler *FileHandler
	folders []string
}

func NewWatcher(logger *zap.SugaredLogger, token string) (*Watcher, error) {
	if w, err := fsnotify.NewWatcher(); err != nil {
		return nil, err
	} else {
		watcher := &Watcher{spy: w, logger: logger.Named("Watcher"), folders: make([]string, 0)}
		watcher.handler = NewFileHandler(watcher.logger, watcher.spy, token)
		return watcher, nil
	}
}

func (w *Watcher) Start() {
	go w.handleEvents()
}

func (w *Watcher) handleEvents() {
	const methodName = "handleEvents"
	defer w.spy.Close()
	for {
		select {
		case event, ok := <-w.spy.Events:
			if !ok {
				w.logger.Errorw(methodName, "Event Spy Ok?", ok)
				return
			}
			if ext, ok := w.handler.HasValidExtension(event.Name); ok {
				w.logger.Infow(methodName, "Extension", ext)
			}
			w.handler.HandleEvent(event, w.folders)
		case err, ok := <-w.spy.Errors:
			w.logger.Infow(methodName, "Case", "Err", "Value", err, "Ok?", ok)
			if !ok {
				w.logger.Errorw(methodName, "Err Spy Ok?", ok)
				return
			}
			w.logger.Errorw(methodName, "Channel Error", err)
		}
	}
}

func (w *Watcher) AddDirectory(path string) error {
	const methodName = "AddDirectory"
	var dirName string

	// need to go through all the directory entries and go from there

	w.logger.Infow(methodName, "Adding Directory", path)
	w.folders = append(w.folders, path)

	dirNames := make([]string, 1, 1)
	dirNames[0] = path

	for {
		if len(dirNames) == 0 {
			break
		}

		dirName, dirNames = dirNames[0], dirNames[1:]

		if watchErr := w.spy.Add(dirName); watchErr != nil {
			return watchErr
		}
		w.logger.Infow(methodName, "Watched Directory:", dirName)
		if de, err := os.ReadDir(dirName); err != nil {
			w.logger.Errorw(methodName, "Directory Read Error", err)
			return err
		} else {
			for _, v := range de {
				if v.IsDir() && v.Name() != ".git" {
					lastDir := len(dirNames)
					dirNames = append(dirNames, dirName+"/"+v.Name())
					w.logger.Infow(methodName, "New Directory", dirNames[lastDir])
				}
			}
		}
	}
	return nil
}
