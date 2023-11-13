package limiter

import "context"

type Limiter interface {
	// Limit trigger rate limit?
	// return true，if trigger rate limit
	Limit(ctx context.Context, key string) (bool, error)
}
