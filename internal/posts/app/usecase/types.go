package usecase

import "time"

type FindPostParams struct {
	PageSize   int64
	PageNumber int64
}

type Post struct {
	postId   int64
	content  string
	postedOn time.Time
	postedBy string
	score    int
}
