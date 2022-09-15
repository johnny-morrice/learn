CREATE TABLE blog_tags (
    id SERIAL PRIMARY KEY,
    blog_id INT REFERENCES blogposts(id),
    tag_id INT REFERENCES tags(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);