package usecase

import "time"

type FindPostParams struct {
	PageSize   int64 `form:"page_size"`
	PageNumber int64 `form:"page_number"`
}

type PostRateParams struct {
	PostId  int64  `json:"post_id"`
	Account string `json:"account"`
	Rate    int32  `json:"rate"`
}

type Post struct {
	PostId   int64     `json:"post_id"`
	Content  string    `json:"content"`
	PostedOn time.Time `json:"posted_on"`
	PostedBy string    `json:"posted_by"`
	Score    int       `json:"score"`
}
