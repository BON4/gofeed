CREATE TABLE Posts (
       post_id bigserial not null unique PRIMARY KEY,
       content text not null,
       posted_on timestamptz NOT NULL,
       posted_by text not null,
       score integer not null
)
