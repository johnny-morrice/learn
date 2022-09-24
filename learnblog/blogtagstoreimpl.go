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

type PostTag struct {
	gorm.Model
	TagID      int
	Tag        Tag
	BlogPostID int
	BlogPost   BlogPost
}

func (store BlogTagStoreImpl) UpdatePostTags(ctx context.Context, blogPostID uint, tags []string) error {
	panic("not implemented")
	// 1. find tags that already exist
	// 2. create tags that dont exist
	// 3. find posttags that already exist
	// 4. create the posttags that dont already exist
	// 5. delete posttags that are no longer required
	// Note: I think it is OK to keep tags that are not linked to posttags anymore
}
