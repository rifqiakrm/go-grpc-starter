package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"grpc-starter/common/config"
)

// FileResponse is the payload for sending email
type FileResponse struct {
	Code    int64        `json:"code"`
	Message string       `json:"message"`
	Data    *FileService `json:"data"`
}

// FileService is the payload for sending email
type FileService struct {
	ImageHost string `json:"image_host"`
	FileName  string `json:"file_name"`
	FileMime  string `json:"file_mime"`
	Folder    string `json:"folder"`
}

// UploadFile upload file to file service
func UploadFile(cfg config.Config, fileName string, pathFolder string) (string, error) {
	var response FileResponse

	filePath := filepath.Clean(fileName)
	pdfFile, err := os.Open(filePath)
	if err != nil {
		log.Println("ERROR OPEN FILE:")
		return "", err
	}
	defer func() {
		err = pdfFile.Close()
	}()

	url := cfg.CloudStorage.FileServiceBaseURL + "/v1/upload"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("folder", pathFolder)
	_ = writer.WriteField("file", fileName)
	wPdfFile, err := writer.CreateFormFile("file", pdfFile.Name())
	if err != nil {
		log.Println("ERROR CREATE FORM:")
		return "", err
	}
	_, err = io.Copy(wPdfFile, pdfFile)
	if err != nil {
		log.Println("ERROR COPY:")
		return "", err
	}
	err = writer.Close()
	if err != nil {
		log.Println("ERROR CLOSE WRITER:")
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, url, payload) //nolint
	if err != nil {
		log.Println("ERROR CREATE NEW REQUEST:")
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR POST:", err)
		return "", err
	}
	defer func() {
		err = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&response)
		if err != nil {
			log.Println("ERROR DECODE BODY RESPONSE:")
			return "", err
		}

		fileURL := response.Data.ImageHost + "/" + response.Data.Folder + "/" + response.Data.FileName
		return fileURL, nil
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("ERROR DECODE ERROR:")
		return "", err
	}
	log.Println("error from file service : ", response)
	err = errors.New("FILE SERVICE ERROR: " + resp.Status)
	return "", err
}
