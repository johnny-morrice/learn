CREATE TABLE blogposts (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL
);