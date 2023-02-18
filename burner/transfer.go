package burner

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type FileUpload struct {
	logger *zap.SugaredLogger
	client *http.Client
}

type FileData struct {
	Filename string `json:"filename"`
	Code     string `json:"code,omitempty"`
}

func NewFileData(filename string, strippedFilename string) (*FileData, error) {
	fd := &FileData{Filename: strippedFilename}
	if data, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else {
		fd.Code = base64.URLEncoding.EncodeToString(data)
	}

	return fd, nil
}

func NewFileUpload(logger *zap.SugaredLogger) *FileUpload {
	return &FileUpload{
		logger: logger,
		client: &http.Client{},
	}
}

func (fu *FileUpload) DeleteFile(strippedFilename string) bool {
	const methodName = "UploadFile"
	fd := &FileData{Filename: strippedFilename}

	var dataReader *bytes.Buffer

	if data, err := json.Marshal(fd); err == nil {
		dataReader = bytes.NewBuffer(data)
		fu.logger.Infow(methodName, "JSON", string(data))
	} else {
		fu.logger.Errorw(methodName, "Marshall Error", err)
		return false
	}

	var response *http.Response
	var req *http.Request
	//if r, err := fu.client.Post("http://localhost:9990/", "application/json", dataReader); err != nil {
	if r, err := http.NewRequest(http.MethodDelete, "http://localhost:9990", dataReader); err != nil {
		fu.logger.Errorw(methodName, "Request Create Error", err)
		return false
	} else {
		req = r
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer naPKnS90HLgXL7CBDIa/NyI6fbzsgLGHVGwz6op9L3hmD57397D/DhidOd6sE1F5")
	if r, err := fu.client.Do(req); err != nil {
		fu.logger.Errorw(methodName, "Delete Error", err)
		return false
	} else {
		response = r
	}
	defer response.Body.Close()

	var body string

	if bodyData, err := io.ReadAll(response.Body); err == nil {
		body = string(bodyData)
	} else {
		fu.logger.Errorw(methodName, "Response Body Error", err)
	}

	fu.logger.Infow(methodName, "Response Code", response.StatusCode, "Response Status", response.Status, "Response Body", body)

	return true
}

func (fu *FileUpload) UploadFile(filename string, strippedFilename string) bool {
	const methodName = "UploadFile"
	var fd *FileData
	if d, err := NewFileData(filename, strippedFilename); err != nil {
		fu.logger.Errorw(methodName, "File Data Error", err)
	} else {
		fd = d
		fu.logger.Infow(methodName, "File Data", fd)
	}

	var dataReader *bytes.Buffer

	if data, err := json.Marshal(fd); err == nil {
		dataReader = bytes.NewBuffer(data)
		fu.logger.Infow(methodName, "JSON", string(data))
	} else {
		fu.logger.Errorw(methodName, "Marshall Error", err)
		return false
	}

	var response *http.Response
	var req *http.Request
	//if r, err := fu.client.Post("http://localhost:9990/", "application/json", dataReader); err != nil {
	if r, err := http.NewRequest(http.MethodPost, "http://localhost:9990", dataReader); err != nil {
		fu.logger.Errorw(methodName, "Request Create Error", err)
		return false
	} else {
		req = r
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer naPKnS90HLgXL7CBDIa/NyI6fbzsgLGHVGwz6op9L3hmD57397D/DhidOd6sE1F5")
	if r, err := fu.client.Do(req); err != nil {
		fu.logger.Errorw(methodName, "Post Error", err)
		return false
	} else {
		response = r
	}
	defer response.Body.Close()

	var body string

	if bodyData, err := io.ReadAll(response.Body); err == nil {
		body = string(bodyData)
	} else {
		fu.logger.Errorw(methodName, "Response Body Error", err)
	}

	fu.logger.Infow(methodName, "Response Code", response.StatusCode, "Response Status", response.Status, "Response Body", body)

	return true
}
