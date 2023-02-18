package burner

import (
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type Watcher struct {
	spy     *fsnotify.Watcher
	logger  *zap.SugaredLogger
	handler *FileHandler
	folders []string
}

func NewWatcher(logger *zap.SugaredLogger) (*Watcher, error) {
	if w, err := fsnotify.NewWatcher(); err != nil {
		return nil, err
	} else {
		watcher := &Watcher{spy: w, logger: logger.Named("Watcher"), folders: make([]string, 0)}
		watcher.handler = NewFileHandler(watcher.logger, watcher.spy)
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
	w.folders = append(w.folders, path)
	return w.spy.Add(path)
}
