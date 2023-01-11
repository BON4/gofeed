CREATE TABLE Posts (
       post_id bigserial not null unique PRIMARY KEY,
       content text not null,
       posted_on timestamptz NOT NULL,
       posted_by text not null,
       score integer not null
);

-- TODO: create index for account
CREATE TABLE RatedPosts (
       post_id bigint not null,
       account text not null,
       rated_score integer not null,
       PRIMARY KEY (post_id, account)
);

ALTER TABLE RatedPosts
      ADD FOREIGN KEY (post_id) REFERENCES Posts (post_id);
