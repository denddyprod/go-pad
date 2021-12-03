package models

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type PadModel struct {
	rdb *redis.Client
}

func NewPadModel(rdb *redis.Client) *PadModel {
	return &PadModel{rdb: rdb}
}

func (pm *PadModel) SavePad(urlPath string, text string) error {
	ctx := context.Background()

	err := pm.rdb.Set(ctx, urlPath, text, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (pm *PadModel) GetPad(urlPath string) (string, error) {
	ctx := context.Background()

	val, err := pm.rdb.Get(ctx, urlPath).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
