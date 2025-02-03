package util

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func SavePicture(file multipart.File, header multipart.FileHeader) (string, error) {
	ext := filepath.Ext(header.Filename)
	newFileName := uuid.New().String() + ext
	savePath := filepath.Join("./dist/upload", newFileName)

	// ensure upload folder exists
	err := os.MkdirAll("./dist/upload", os.ModePerm)
	if err != nil {
		return "", err
	}

	// distnation file
	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}

	// copt picture to new file
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	saveUrl := "/upload/" + newFileName

	return saveUrl, nil
}
