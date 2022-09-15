CREATE TABLE blog_tags (
    id SERIAL PRIMARY KEY,
    blog_id INT REFERENCES blogposts(id),
    tag_id INT REFERENCES tags(id)
);