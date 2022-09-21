package main

import (
	"context"

	"gorm.io/gorm"
)

type BlogTagStoreImpl struct {
	DB *gorm.DB
}

type Tag struct {
	gorm.Model
	Name string
}

type BlogTag struct {
	gorm.Model
	TagID      int
	Tag        Tag
	BlogPostID int
	BlogPost   BlogPost
}

func (store BlogTagStoreImpl) AddPostTags(ctx context.Context, blogPostID uint, tags []string) error {
	panic("not implemented")
}
