package burner

import (
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

var (
	fileExtensions = []string{".js", ".script"}
)

type FileHandler struct {
	logger   *zap.SugaredLogger
	uploader *FileUpload
	spy      *fsnotify.Watcher
}

func NewFileHandler(l *zap.SugaredLogger, s *fsnotify.Watcher) *FileHandler {
	return &FileHandler{logger: l.Named("FileHandler"), uploader: NewFileUpload(l), spy: s}
}

func (fh *FileHandler) HasValidExtension(filename string) (string, bool) {
	for _, v := range fileExtensions {
		if strings.HasSuffix(filename, v) {
			return v, true
		}
	}
	return "", false
}

func (fh *FileHandler) HasSubdir(filename string) bool {
	return strings.Contains(filename, "/")
}

func (fh *FileHandler) stripFilename(filename string, watchedFolders []string) string {
	fn := filename
	for _, v := range watchedFolders {
		if strings.HasPrefix(filename, v) {
			fh.logger.Debugw("stripFilename", "Filename", filename, "Folder", v)
			fn = fn[(len(v) + 1):]
		}
	}
	if fh.HasSubdir(fn) {
		fn = "/" + fn
	}
	return fn
}

func (fh *FileHandler) createDirectory(filename string) bool {
	const methodName = "createDirectory"
	if finfo, err := os.Stat(filename); err != nil {
		fh.logger.Errorw(methodName, "Stat File", filename, "Error", err)
		return false
	} else if !finfo.IsDir() {
		fh.logger.Infow(methodName, "Stat File", filename, "Type", "File")
		return false
	}
	fh.spy.Add(filename)
	return true
}

func (fh *FileHandler) createFile(filename, strippedFilename string) {
	fh.uploader.UploadFile(filename, strippedFilename)
}

func (fh *FileHandler) handleCreate(evt fsnotify.Event, watchedFolders []string) {
	sfn := fh.stripFilename(evt.Name, watchedFolders)
	fh.logger.Infow("Creating", "File", evt.Name, "Stripped", sfn, "Has Subdir?", fh.HasSubdir(sfn))

	if fh.createDirectory(evt.Name) {
		return
	}

	if _, ok := fh.HasValidExtension(sfn); ok {
		fh.createFile(evt.Name, sfn)
	}
}

func (fh *FileHandler) handleDelete(evt fsnotify.Event, watchedFolders []string) {
	sfn := fh.stripFilename(evt.Name, watchedFolders)
	fh.logger.Infow("Deleting", "File", evt.Name, "Stripped", sfn, "Has Subdir?", fh.HasSubdir(sfn))
	fh.uploader.DeleteFile(sfn)
}

func (fh *FileHandler) HandleEvent(evt fsnotify.Event, watchedFolders []string) {
	const methodName = "HandleEvent"
	switch {
	case evt.Op.Has(fsnotify.Create):
		fh.handleCreate(evt, watchedFolders)
	case evt.Op.Has(fsnotify.Chmod):
		fh.logger.Infow("CHMOD")
		fh.uploader.UploadFile(evt.Name, fh.stripFilename(evt.Name, watchedFolders))
	case evt.Op.Has(fsnotify.Remove):
		fh.logger.Infow("Deleting")
		fh.handleDelete(evt, watchedFolders)
	case evt.Op.Has(fsnotify.Write):
		fh.logger.Infow("Writing")
	default:
		fh.logger.Infow("Other")
	}
	/*
		strippedName := evt.Name
		fh.logger.Debugw(methodName, "Full Filename", evt.Name, "Op", evt.Op.String())
		if _, ok := fh.HasValidExtension(evt.Name); !ok {
			fh.logger.Infow(methodName, "Valid Extension?", false)
			return
		}

		if finfo, err := os.Stat(evt.Name); err != nil {
			fh.logger.Errorw(methodName, "Stat File", evt.Name, "Error", err)
			return
		} else if finfo.IsDir() {
			fh.logger.Infow(methodName, "Stat File", evt.Name, "Type", "Directory")
			return
		}

		for _, v := range watchedFolders {
			if strings.HasPrefix(strippedName, v) {
				strippedName = strippedName[(len(v) + 1):]
				break
			}
		}
		fh.logger.Debugw(methodName, "Stripped Filename", strippedName, "Op", evt.Op.String(), "Has Subdir?", fh.HasSubdir(strippedName))
		if fh.HasSubdir(strippedName) {
			strippedName = "/" + strippedName
		}

		if evt.Op.Has(fsnotify.Remove) {
			fh.logger.Infow(methodName, "Removing File", evt.Name)
		}
		if evt.Op.Has(fsnotify.Write) || evt.Op.Has(fsnotify.Create) {
			fh.logger.Infow(methodName, "Posting File", evt.Name)
			fh.uploader.UploadFile(evt.Name, strippedName)
		}*/
}
