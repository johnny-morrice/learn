package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type BlogPostStoreImpl struct {
	DB *gorm.DB
}

type BlogPost struct {
	gorm.Model
	UUID  string
	Title string
	Body  string
}

func (store BlogPostStoreImpl) AddPost(ctx context.Context, post BlogPost) error {
	err := store.DB.WithContext(ctx).Create(&post).Error
	if err != nil {
		return fmt.Errorf("failed to create blogpost: %w", err)
	}

	return nil
}
func (store BlogPostStoreImpl) CountPosts(ctx context.Context) (int64, error) {
	count := int64(0)
	err := store.DB.Model(&BlogPost{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}
	return count, nil
}
func (store BlogPostStoreImpl) GetPost(ctx context.Context, postID string) (*BlogPost, error) {
	post := &BlogPost{}
	err := store.DB.WithContext(ctx).First(post, "uuid = ?", postID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get blogpost: %w", err)
	}
	return post, nil
}
func (store BlogPostStoreImpl) GetPostsPage(ctx context.Context, offset, limit int) ([]BlogPost, error) {
	result := []BlogPost{}
	err := store.DB.WithContext(ctx).Offset(offset).Limit(limit).Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get blogpost page: %w", err)
	}
	return result, nil
}
