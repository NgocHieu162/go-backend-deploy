package storage_impl

import (
	"context"
	"fmt"
	"go-backend/internal/common/env"
	"go-backend/internal/common/storage"
	"go-backend/internal/dto"
	"log"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryStorage struct {
	cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryStorage(env *env.Env) storage.FileStorage {
	cld, err := cloudinary.NewFromURL(env.CloudinaryURL)

	if err != nil {
		log.Fatal(err.Error())
	}

	return &cloudinaryStorage{
		cloudinary: cld,
	}
}

// Upload implements [storage.FileStorage].
func (c *cloudinaryStorage) Upload(input dto.UploadInput) (*dto.UploadReturn, error) {
	file, err := input.FileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileName := fmt.Sprintf("%s%d", input.Prefix, time.Now().UnixNano())

	uploadResult, err := c.cloudinary.Upload.Upload(context.Background(), file, uploader.UploadParams{
		PublicID: fileName,
		Folder:   input.Folder,
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\n %+v \n\n", uploadResult)

	return &dto.UploadReturn{
		FileName: fileName,
		FullPath: uploadResult.PublicID,
		Url:      uploadResult.SecureURL,
	}, nil
}

// Delete implements [storage.FileStorage].
func (c *cloudinaryStorage) Delete(fileName string) error {
	_, err := c.cloudinary.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: fileName,
	})
	if err != nil{
		return err
	}

	return nil
}
