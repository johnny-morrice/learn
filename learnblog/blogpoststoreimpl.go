package main

import (
	"context"

	"gorm.io/gorm"
)

type BlogPostStoreImpl struct {
	DB *gorm.DB
}

func (store BlogPostStoreImpl) AddPost(ctx context.Context, post BlogPostRecord) error {
	panic("not implemented")
}
func (store BlogPostStoreImpl) CountPosts(ctx context.Context) (int, error) {
	panic("not implemented")
}
func (store BlogPostStoreImpl) GetPost(ctx context.Context, postID string) (*BlogPostRecord, error) {
	panic("not implemented")
}
func (store BlogPostStoreImpl) GetPostsPage(ctx context.Context, offset, limit int) ([]BlogPostRecord, error) {
	panic("not implemented")
}
