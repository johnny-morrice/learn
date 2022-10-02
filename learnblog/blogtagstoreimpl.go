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

func (store BlogTagStoreImpl) FindTagsByNames(ctx context.Context, tags []string) ([]Tag, error) {
	out := []Tag{}
	err := store.DB.Find(&out, "name IN ?", tags).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (store BlogTagStoreImpl) CreateTags(ctx context.Context, tagNames []string) ([]Tag, error) {
	tags := []Tag{}

	for _, t := range tagNames {
		tags = append(tags, Tag{Name: t})
	}

	err := store.DB.Create(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func DiffTagNames(tags []Tag, tagNames []string) []string {
	diff := []string{}
	existing := map[string]struct{}{}

	for _, t := range tags {
		existing[t.Name] = struct{}{}
	}
	for _, name := range tagNames {
		_, present := existing[name]
		if !present {
			diff = append(diff, name)
		}
	}

	return diff
}

func (store BlogTagStoreImpl) UnlinkedTagsForPost(ctx context.Context, blogPostID uint, tags []string) ([]Tag, error) {
	out := []Tag{}

	const sql = "SELECT * FROM tags where name IN ? AND id NOT IN (SELECT tag_id FROM post_tags where blog_post_id = ?)"

	err := store.DB.Raw(sql, tags, blogPostID).Scan(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (store BlogTagStoreImpl) CreatePostTags(ctx context.Context, blogPostID uint, tags []Tag) ([]PostTag, error) {
	postTags := []PostTag{}

	for _, t := range tags {
		postTags = append(postTags, PostTag{BlogPostID: int(blogPostID), TagID: int(t.ID)})
	}

	err := store.DB.Create(&postTags).Error
	if err != nil {
		return nil, err
	}

	return postTags, nil
}

func (store BlogTagStoreImpl) DeletePostTagsNotInNames(ctx context.Context, blogPostID uint, tagNames []string) error {
	const sql = "DELETE FROM tags WHERE id IN (SELECT pt.tag_id FROM post_tags pt JOIN tags t ON pt.tag_id = t.id WHERE pt.blog_post_id = ? AND t.name NOT IN ?)"

	return store.DB.Exec(sql, blogPostID, tagNames).Error
}

func (store BlogTagStoreImpl) UpdatePostTags(ctx context.Context, blogPostID uint, tagNames []string) error {
	// 1. find tags that already exist
	// 2. create tags that dont exist
	// 3. find posttags that already exist
	// 4. create the posttags that dont already exist
	// 5. delete posttags that are no longer required
	// Note: I think it is OK to keep tags that are not linked to posttags anymore

	existingTags, err := store.FindTagsByNames(ctx, tagNames)
	if err != nil {
		return err
	}

	newNames := DiffTagNames(existingTags, tagNames)

	_, err = store.CreateTags(ctx, newNames)
	if err != nil {
		return err
	}

	unlinkedTags, err := store.UnlinkedTagsForPost(ctx, blogPostID, tagNames)
	if err != nil {
		return err
	}

	_, err = store.CreatePostTags(ctx, blogPostID, unlinkedTags)
	if err != nil {
		return err
	}

	return store.DeletePostTagsNotInNames(ctx, blogPostID, tagNames)
}
