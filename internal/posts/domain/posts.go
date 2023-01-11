package domain

import (
	"errors"
	"fmt"
	"time"
)

type Post struct {
	postId   int64 //set on database save. If not set equals to -1.
	content  string
	postedOn time.Time //set on init now
	postedBy string
	score    int32 //init with 0
}

func (p *Post) ID() int64 {
	return p.postId
}

func (p *Post) UpWote() {
	p.score++
}

func (p *Post) IsHiden() bool {
	return p.score < 0
}

func (p *Post) DownWote() {
	p.score--
}

func (p *Post) Content() string {
	return p.content
}

func (p *Post) PostedOn() time.Time {
	return p.postedOn
}

func (p *Post) PostedBy() string {
	return p.postedBy
}

func (p *Post) Score() int32 {
	return p.score
}

type TooLongContent struct {
	MaxContentLength      int
	ProvidedContentLength int
}

func (e TooLongContent) Error() string {
	return fmt.Sprintf(
		"Provided content is too long, maximum length: %d, provided content length: %d",
		e.MaxContentLength,
		e.ProvidedContentLength,
	)
}

// UnmarshalPostFromDatabase - unmarshals psot from the database.
//
// It should be used only for unmarshalling from the database!
func UnmarshalPostFromDatabase(
	postId int64,
	content string,
	postedOn time.Time,
	poster string,
	score int32) (*Post, error) {

	p, err := NewPost(poster, content)
	if err != nil {
		return nil, err
	}

	// Setting field that not set by default
	p.postId = postId
	p.postedOn = postedOn
	p.score = score

	return p, nil
}

func NewPost(poster, content string) (*Post, error) {
	if len(poster) == 0 {
		return nil, errors.New("Username must be not empty")
	}

	if len(content) > 400 {
		return nil, TooLongContent{
			MaxContentLength:      400,
			ProvidedContentLength: len(content),
		}
	}

	return &Post{
		postId:   -1,
		content:  content,
		postedBy: poster,
		postedOn: time.Now(),
	}, nil
}
