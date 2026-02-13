package storage

import (
	"context"
	"strings"
	"time"

	"pet-services-api/internal/entities"
)

const PhotoSignedURLTTL = 15 * time.Minute

func SignUserPhotos(ctx context.Context, storageService ObjectStorage, user *entities.User) error {
	if user == nil || storageService == nil || len(user.Photos) == 0 {
		return nil
	}

	for i := range user.Photos {
		key := strings.TrimSpace(user.Photos[i].URL)
		if key == "" {
			continue
		}
		if strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
			continue
		}
		if !strings.Contains(key, "/") {
			key = "users/" + user.ID + "/" + key
		}

		url, err := storageService.GenerateReadURL(ctx, key, PhotoSignedURLTTL)
		if err != nil {
			return err
		}
		user.Photos[i].URL = url
	}

	return nil
}

func SignUsersPhotos(ctx context.Context, storageService ObjectStorage, users []*entities.User) error {
	if storageService == nil || len(users) == 0 {
		return nil
	}
	for _, user := range users {
		if err := SignUserPhotos(ctx, storageService, user); err != nil {
			return err
		}
	}
	return nil
}
