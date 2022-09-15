package main

import (
	"context"
	"errors"
)

type BlogPostStoreMemoryImpl struct {
	Posts []BlogPostRecord
}

func (store *BlogPostStoreMemoryImpl) AddPost(ctx context.Context, post BlogPostRecord) error {
	if ctx.Err() != nil {
		return errors.New("context expired")
	}
	post.ID = uint(len(store.Posts))
	store.Posts = append(store.Posts, post)
	return nil
}

func (store *BlogPostStoreMemoryImpl) CountPosts(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, errors.New("context expired")
	}
	return len(store.Posts), nil
}

func (store *BlogPostStoreMemoryImpl) GetPost(ctx context.Context, postID string) (*BlogPostRecord, error) {
	if ctx.Err() != nil {
		return nil, errors.New("context expired")
	}
	for _, rec := range store.Posts {
		if postID == rec.UUID {
			return &rec, nil
		}
	}
	return nil, errors.New("could not find post")
}

func (store *BlogPostStoreMemoryImpl) GetPostsPage(ctx context.Context, offset, limit int) ([]BlogPostRecord, error) {
	if ctx.Err() != nil {
		return nil, errors.New("context expired")
	}

	if offset < 0 || limit < 1 {
		return nil, errors.New("invalid parameters for GetPostsPage")
	}

	if offset >= len(store.Posts) {
		return []BlogPostRecord{}, nil
	}

	if offset+limit > len(store.Posts) {
		limit = len(store.Posts) - offset
	}

	return store.Posts[offset : offset+limit], nil
}
