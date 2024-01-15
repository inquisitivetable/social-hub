package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type IImageService interface {
	SaveImage(file multipart.File, fileHeader multipart.FileHeader) (string, error)
}

func SaveImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Get the file extension
	fileExtension := strings.Split(fileHeader.Filename, ".")[1]

	// Generate new file name
	newFileName := fmt.Sprintf("%s.%s", uuid.NewV4().String(), fileExtension)

	imagePath := "images"

	// Create folder if not exists
	err := os.MkdirAll(imagePath, os.ModePerm)
	if err != nil {
		log.Println("Error with creating folder", err)
		return "", err
	}

	// Create new file
	newFile, err := os.Create(fmt.Sprintf("%s/%s", imagePath, newFileName))
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer newFile.Close()

	// Copy the uploaded file to the created file
	if _, err := io.Copy(newFile, file); err != nil {
		log.Println(err)
		return "", err
	}

	return newFileName, nil
}
