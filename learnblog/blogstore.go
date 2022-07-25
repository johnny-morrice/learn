package main

import (
	"errors"
)

type BlogPostRecord struct {
	ID    int
	UUID  string
	Title string
	Body  string
}

type BlogStore struct {
	Posts []BlogPostRecord
}

func (store *BlogStore) AddPost(post BlogPostRecord) error {
	post.ID = len(store.Posts)
	store.Posts = append(store.Posts, post)
	return nil
}

func (store *BlogStore) CountPosts() (int, error) {
	return len(store.Posts), nil
}

func (store *BlogStore) GetPost(postID string) (*BlogPostRecord, error) {
	for _, rec := range store.Posts {
		if postID == rec.UUID {
			return &rec, nil
		}
	}
	return nil, errors.New("could not find post")
}

func (store *BlogStore) GetPostsPage(offset, limit int) ([]BlogPostRecord, error) {
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
