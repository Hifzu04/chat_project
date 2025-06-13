package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	cloudName := os.Getenv("CLOUDNIARY_CLOUDNAME")
	apiKey := os.Getenv("CLOUDNIARY_APIKEYS")
	apiSecret := os.Getenv("CLOUDNIARY_APISECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", fmt.Errorf("cloudinary init error: %v", err)
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder:   "chatapp",
		PublicID: fileHeader.Filename,
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary upload error: %v", err)
	}

	return uploadResult.SecureURL, nil
}
