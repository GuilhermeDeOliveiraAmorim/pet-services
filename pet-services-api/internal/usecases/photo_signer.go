package usecases

import (
	"context"
	"strings"
	"time"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/storage"
)

const photoSignedURLTTL = 15 * time.Minute

func signUserPhotos(ctx context.Context, storageService storage.ObjectStorage, user *entities.User) error {
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

		url, err := storageService.GenerateReadURL(ctx, key, photoSignedURLTTL)
		if err != nil {
			return err
		}
		user.Photos[i].URL = url
	}

	return nil
}

func signUsersPhotos(ctx context.Context, storageService storage.ObjectStorage, users []*entities.User) error {
	if storageService == nil || len(users) == 0 {
		return nil
	}
	for _, user := range users {
		if err := signUserPhotos(ctx, storageService, user); err != nil {
			return err
		}
	}
	return nil
}
