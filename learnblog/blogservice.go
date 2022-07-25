package main

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type BlogPost struct {
	UUID  string
	Title string
	Body  string
}

func (post BlogPost) Validate() error {
	if post.Title == "" {
		return errors.New("expected post to have title")
	}
	_, err := uuid.Parse(post.UUID)
	if err != nil {
		return errors.Wrap(err, "failed to parse post UUID")
	}
	return nil
}

type BlogService struct {
	Store *BlogStore
}

type BlogPostPage struct {
	Total  int
	Limit  int
	Offset int
	Posts  []BlogPost
}

func BlogRecToBlogPost(rec BlogPostRecord) BlogPost {
	return BlogPost{
		UUID:  rec.UUID,
		Title: rec.Title,
		Body:  rec.Body,
	}
}

func BlogPostToBlogRec(post BlogPost) BlogPostRecord {
	return BlogPostRecord{
		UUID:  post.UUID,
		Title: post.Title,
		Body:  post.Body,
	}
}

func (srv BlogService) GetPost(postID string) (*BlogPost, error) {
	rec, err := srv.Store.GetPost(postID)
	if err != nil {
		return nil, err
	}
	post := BlogRecToBlogPost(*rec)
	return &post, nil
}

func (srv BlogService) AddPost(post BlogPost) error {
	record := BlogPostRecord{
		UUID:  post.UUID,
		Title: post.Title,
		Body:  post.Body,
	}
	return srv.Store.AddPost(record)
}

func (srv BlogService) GetPostsPage(offset, limit int) (*BlogPostPage, error) {
	if offset < 0 || limit < 1 {
		return nil, errors.New("invalid parameters for GetPostsPage")
	}

	postRecords, err := srv.Store.GetPostsPage(offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := srv.Store.CountPosts()

	if err != nil {
		return nil, err
	}

	posts := []BlogPost{}

	for _, rec := range postRecords {
		posts = append(posts, BlogPost{
			UUID:  rec.UUID,
			Title: rec.Title,
			Body:  rec.Body,
		})
	}

	page := &BlogPostPage{
		Total:  total,
		Offset: offset,
		Limit:  limit,
		Posts:  posts,
	}

	return page, nil
}
