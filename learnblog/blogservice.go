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
	Tags  []string
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

type BlogPostStore interface {
	AddPost(ctx context.Context, post *BlogPost) error
	CountPosts(ctx context.Context) (int64, error)
	CountPostsWithTag(ctx context.Context, tag string) (int64, error)
	GetPost(ctx context.Context, postID string) (*BlogPost, error)
	GetPostsPage(ctx context.Context, offset, limit int) ([]BlogPost, error)
	GetPostsPageByTag(ctx context.Context, offset, limit int, tag string) ([]BlogPost, error)
}

type BlogTagStore interface {
	UpdatePostTags(ctx context.Context, blogPostID uint, tags []string) error
}

type BlogService struct {
	PostStore BlogPostStore
	TagStore  BlogTagStore
}

type BlogPostPage struct {
	Total  int64
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
	rec, err := srv.PostStore.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	post := BlogRecToBlogPost(*rec)
	return &post, nil
}

func (srv BlogService) AddPost(ctx context.Context, post BlogPostViewModel) error {
	record := &BlogPost{
		UUID:  post.UUID,
		Title: post.Title,
		Body:  post.Body,
	}
	err := srv.PostStore.AddPost(ctx, record)
	if err != nil {
		return err
	}
	return srv.TagStore.UpdatePostTags(ctx, record.ID, post.Tags)
}

func (srv BlogService) GetPostsPage(ctx context.Context, offset, limit int) (*BlogPostPage, error) {
	if offset < 0 || limit < 1 {
		return nil, errors.New("invalid parameters for GetPostsPage")
	}

	postRecords, err := srv.PostStore.GetPostsPage(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := srv.PostStore.CountPosts(ctx)

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

func (srv BlogService) GetPostsPageByTag(ctx context.Context, offset, limit int, tag string) (*BlogPostPage, error) {
	if offset < 0 || limit < 1 {
		return nil, errors.New("invalid parameters for GetPostsPage")
	}

	postRecords, err := srv.PostStore.GetPostsPageByTag(ctx, offset, limit, tag)
	if err != nil {
		return nil, err
	}

	total, err := srv.PostStore.CountPostsWithTag(ctx, tag)

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
