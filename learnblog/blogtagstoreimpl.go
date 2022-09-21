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

// below should be called PostTag
type BlogTag struct {
	gorm.Model
	TagID      int
	Tag        Tag
	BlogPostID int
	BlogPost   BlogPost
}

func (store BlogTagStoreImpl) AddPostTags(ctx context.Context, blogPostID uint, tags []string) error {
	panic("not implemented")
	// 1. find tags that already exist
	// 2. create tags that dont exist
	// 3. find blogtags that already exist
	// 4. create the blogtags that dont already exist
}
