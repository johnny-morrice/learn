package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type BlogPostStoreImpl struct {
	DB *gorm.DB
}

type BlogPostRecord struct {
	gorm.Model
	UUID  string
	Title string
	Body  string
}

func (store BlogPostStoreImpl) AddPost(ctx context.Context, post BlogPostRecord) error {
	err := store.DB.WithContext(ctx).Create(&post).Error
	if err != nil {
		return fmt.Errorf("failed to create blogpost: %w", err)
	}

	return nil
}
func (store BlogPostStoreImpl) CountPosts(ctx context.Context) (int, error) {
	panic("not implemented")
}
func (store BlogPostStoreImpl) GetPost(ctx context.Context, postID string) (*BlogPostRecord, error) {
	post := &BlogPostRecord{}
	err := store.DB.WithContext(ctx).First(post, "uuid = ?", postID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get blogpost: %w", err)
	}
	return post, nil
}
func (store BlogPostStoreImpl) GetPostsPage(ctx context.Context, offset, limit int) ([]BlogPostRecord, error) {
	panic("not implemented")
}
