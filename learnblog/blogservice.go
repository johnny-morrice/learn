package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type BlogPostViewModel struct {
	UUID  string
	Title string
	Body  string
}

func (post BlogPostViewModel) Validate() error {
	if post.Title == "" {
		return errors.New("expected post to have title")
	}
	_, err := uuid.Parse(post.UUID)
	if err != nil {
		return errors.Wrap(err, "failed to parse post UUID")
	}
	return nil
}

type BlogStore interface {
	AddPost(ctx context.Context, post BlogPost) error
	CountPosts(ctx context.Context) (int, error)
	GetPost(ctx context.Context, postID string) (*BlogPost, error)
	GetPostsPage(ctx context.Context, offset, limit int) ([]BlogPost, error)
}

type BlogService struct {
	Store BlogStore
}

type BlogPostPage struct {
	Total  int
	Limit  int
	Offset int
	Posts  []BlogPostViewModel
}

func BlogRecToBlogPost(rec BlogPost) BlogPostViewModel {
	return BlogPostViewModel{
		UUID:  rec.UUID,
		Title: rec.Title,
		Body:  rec.Body,
	}
}

func BlogPostToBlogRec(post BlogPostViewModel) BlogPost {
	return BlogPost{
		UUID:  post.UUID,
		Title: post.Title,
		Body:  post.Body,
	}
}

func (srv BlogService) GetPost(ctx context.Context, postID string) (*BlogPostViewModel, error) {
	rec, err := srv.Store.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	post := BlogRecToBlogPost(*rec)
	return &post, nil
}

func (srv BlogService) AddPost(ctx context.Context, post BlogPostViewModel) error {
	record := BlogPost{
		UUID:  post.UUID,
		Title: post.Title,
		Body:  post.Body,
	}
	return srv.Store.AddPost(ctx, record)
}

func (srv BlogService) GetPostsPage(ctx context.Context, offset, limit int) (*BlogPostPage, error) {
	if offset < 0 || limit < 1 {
		return nil, errors.New("invalid parameters for GetPostsPage")
	}

	postRecords, err := srv.Store.GetPostsPage(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := srv.Store.CountPosts(ctx)

	if err != nil {
		return nil, err
	}

	posts := []BlogPostViewModel{}

	for _, rec := range postRecords {
		posts = append(posts, BlogPostViewModel{
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
