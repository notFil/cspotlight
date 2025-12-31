package util

import (
	"errors"
	"io"
	"mime/multipart"
	"slices"

	"github.com/gabriel-vasile/mimetype"
)

const ImageMaxSize = 2 << 20

var ImageAllowedExts = []string{"image/jpeg", "image/png", "image/gif"}

func ValidateImage(image *multipart.FileHeader) error {
	if image.Size > ImageMaxSize {
		return errors.New("file size too large")
	}
	f, err := image.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	if err := validateImageMimeType(f); err != nil {
		return err
	}
	return nil
}

func validateImageMimeType(f multipart.File) error {
	mimeType, err := mimetype.DetectReader(f)
	if err != nil {
		return err
	}
	if !slices.Contains(ImageAllowedExts, mimeType.String()) {
		return errors.New("invalid image format")
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	return nil
}
