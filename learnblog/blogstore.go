package main

import (
	"errors"
	"log"
)

type BlogPostRecord struct {
	ID    int
	UUID  string
	Title string
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

func (store *BlogStore) GetPostsPage(offset, limit int) ([]BlogPostRecord, error) {
	log.Printf("getting post with offset, limit: %v, %v", offset, limit)
	if offset < 0 || limit < 1 {
		return nil, errors.New("invalid parameters for GetPostsPage")
	}

	if offset >= len(store.Posts) {
		return []BlogPostRecord{}, nil
	}

	if offset+limit > len(store.Posts) {
		limit = len(store.Posts) - offset
		log.Printf("changed limit to: %v", limit)
	}

	return store.Posts[offset : offset+limit], nil
}
