package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"be-file-uploader/internal/models"

	"github.com/redis/go-redis/v9"
)

func (r *Repository) CreateSession(ctx context.Context, session *models.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	ttl := time.Until(session.ExpiresAt)

	pipe := r.redis.Pipeline()

	pipe.Set(ctx, fmt.Sprintf("session:%s", session.ID), data, ttl)

	pipe.Set(ctx, fmt.Sprintf("session_hash:%s", session.RefreshHash), session.ID, ttl)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *Repository) SearchSessionByRefreshHash(ctx context.Context, hash string) (*models.Session, error) {
	sessionID, err := r.redis.Get(ctx, fmt.Sprintf("session_hash:%s", hash)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return r.SearchSessionByID(ctx, sessionID)
}

func (r *Repository) SearchSessionByID(ctx context.Context, id string) (*models.Session, error) {
	data, err := r.redis.Get(ctx, fmt.Sprintf("session:%s", id)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	session := new(models.Session)
	if err := json.Unmarshal([]byte(data), session); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *Repository) TerminateSession(ctx context.Context, sessionID string) (bool, error) {
	err := r.redis.Del(ctx, fmt.Sprintf("session:%s", sessionID)).Err()
	if err != nil {
		return false, err
	}

	return true, err
}
