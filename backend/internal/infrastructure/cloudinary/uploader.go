package cloudinary

import (
	"context"
	"fmt"

	cldgo "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"featherweight/internal/config"
)

// Uploader wraps the Cloudinary SDK and provides Upload/Delete operations.
type Uploader struct {
	client *cldgo.Cloudinary
	folder string
}

// NewUploader creates a Cloudinary Uploader from config credentials.
func NewUploader(cfg config.Config) (*Uploader, error) {
	client, err := cldgo.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)
	if err != nil {
		return nil, fmt.Errorf("init cloudinary client: %w", err)
	}

	return &Uploader{
		client: client,
		folder: cfg.CloudinaryFolder,
	}, nil
}

// Upload sends a local file to Cloudinary and returns (secureURL, publicID, error).
// publicID is the Cloudinary asset identifier needed for future deletion.
func (u *Uploader) Upload(ctx context.Context, localPath, publicID string) (secureURL string, fullPublicID string, err error) {
	result, err := u.client.Upload.Upload(ctx, localPath, uploader.UploadParams{
		PublicID:     publicID,
		Folder:       u.folder,
		ResourceType: "image",
	})
	if err != nil {
		return "", "", fmt.Errorf("cloudinary upload: %w", err)
	}

	return result.SecureURL, result.PublicID, nil
}

// Delete removes an asset from Cloudinary by its public_id.
func (u *Uploader) Delete(ctx context.Context, publicID string) error {
	_, err := u.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	})
	if err != nil {
		return fmt.Errorf("cloudinary delete %s: %w", publicID, err)
	}

	return nil
}
