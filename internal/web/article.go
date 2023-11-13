package web

import "gorm.io/gorm/logger"

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(l logger.LoggerV1, svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		l:   l,
		svc: svc,
	}
}
