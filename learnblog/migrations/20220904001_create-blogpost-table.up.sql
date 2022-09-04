CREATE TABLE blogposts (
    id SERIAL,
    uuid VARCHAR(16) NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL
)