package domain

import (
	"time"
)

type Article struct {
	Id        int64
	Title     string
	Content   string
	Author    Author
	Status    ArticleStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a Article) Abstract() string {
	str := []rune(a.Content)
	// content abstract
	if len(str) > 128 {
		str = str[:128]
	}
	return string(str)
}

type ArticleStatus uint8

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

const (
	// ArticleStatusUnknown
	ArticleStatusUnknown = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

type Author struct {
	Id   int64
	Name string
}
