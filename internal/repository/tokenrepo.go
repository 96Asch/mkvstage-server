package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/go-redis/redis/v8"
)

type redisTokenRepo struct {
	R *redis.Client
}

func NewRedisTokenRepository(client *redis.Client) *redisTokenRepo {
	return &redisTokenRepo{
		R: client,
	}
}

func (tr redisTokenRepo) GetAll(ctx context.Context, uid int64) (*[]domain.RefreshToken, error) {
	matcher := fmt.Sprintf("%d:*", uid)
	scan := tr.R.Scan(ctx, 0, matcher, 10)
	if err := scan.Err(); err != nil {
		return nil, domain.NewInternalErr()
	}

	iterator := scan.Iterator()
	tokens := make([]domain.RefreshToken, 0)
	for iterator.Next(ctx) {
		split := strings.Split(iterator.Val(), ":")

		if len(split) != 2 {
			return nil, domain.NewInternalErr()
		}

		id, err := strconv.Atoi(split[0])
		if err != nil {
			return nil, domain.NewInternalErr()
		}

		tokens = append(tokens, domain.RefreshToken{
			UserID:  int64(id),
			Refresh: split[1],
		})
	}

	return &tokens, nil
}

func (tr redisTokenRepo) Create(ctx context.Context, token *domain.RefreshToken) error {
	key := fmt.Sprintf("%d:%s", token.UserID, token.Refresh)
	log.Println(key)
	cmd := tr.R.Set(ctx, key, 0, token.ExpirationDuration)
	if err := cmd.Err(); err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (tr redisTokenRepo) Delete(ctx context.Context, uid int64, refresh string) error {
	key := fmt.Sprintf("%d:%s", uid, refresh)
	err := tr.R.Del(ctx, key).Err()
	if err != nil {
		return domain.NewInternalErr()
	}
	return nil
}

func (tr redisTokenRepo) DeleteAll(ctx context.Context, uid int64) error {
	tokens, err := tr.GetAll(ctx, uid)
	if err != nil {
		return domain.NewInternalErr()
	}

	for _, token := range *tokens {
		key := fmt.Sprintf("%d:%s", token.UserID, token.Refresh)
		err := tr.R.Del(ctx, key).Err()
		if err != nil {
			return domain.NewInternalErr()
		}
	}

	return nil
}
