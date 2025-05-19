package repo

import (
	"context"
	"market/internal/model"
)

type CoinRepo interface {
	FindByUnit(ctx context.Context, unit string) (coin *model.Coin, err error)
	FindById(ctx context.Context, id int64) (coin *model.Coin, err error)
	FindAll(ctx context.Context) ([]*model.Coin, error)
}
