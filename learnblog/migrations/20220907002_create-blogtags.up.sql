CREATE TABLE blogtags (
    id SERIAL PRIMARY KEY,
    blog_id INT REFERENCES blogposts(id),
    tag_id INT REFERENCES tags(id)
);