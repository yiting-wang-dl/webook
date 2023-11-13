package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/webook/internal/domain"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, res []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, art domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, res domain.Article) error
}

type ArticleRedisCache struct {
	client redis.Cmdable
}

func NewArticleRedisCache(client redis.Cmdable) ArticleCache {
	return &ArticleRedisCache{
		client: client,
	}
}

func (a *ArticleRedisCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
}

func (a *ArticleRedisCache) SetFirstPage(ctx context.Context, uid int64, res []domain.Article) error {

}

func (a *ArticleRedisCache) DelFirstPage(ctx context.Context, uid int64) error {

}

func (a *ArticleRedisCache) Get(ctx context.Context, id int64) (domain.Article, error) {

}

func (a *ArticleRedisCache) Set(ctx context.Context, art domain.Article) error {

}
func (a *ArticleRedisCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {

}
func (a *ArticleRedisCache) SetPub(ctx context.Context, res domain.Article) error {

}
