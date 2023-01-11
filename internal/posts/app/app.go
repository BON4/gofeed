package app

import "github.com/BON4/gofeed/internal/posts/app/usecase"

type Application struct {
	CreatePost usecase.CreatePostHandler
	DeletePost usecase.DeletePostHandler
	ListPost   usecase.ListPostHandler
	RatePost   usecase.RatePostHandler
}
