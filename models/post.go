package models

import (
	"forum/pkg/validator"
	"strconv"
	"time"
)

type Post struct {
	PostID       int
	UserID       int
	UserName     string
	Title        string
	Content      string
	ImageName    string
	Created      time.Time
	Like         int
	Dislike      int
	Comment      *[]Comment
	Categories   map[int]string
	IsLiked      int
	CommentCount int
}

type Comment struct {
	CommentID int
	PostID    int
	UserID    int
	UserName  string
	Content   string
	Created   time.Time
	Like      string
	Dislike   string
	IsLiked   int
}

type CommentForm struct {
	PostID  int
	UserID  int
	Content string
	Token   string
	validator.Validator
}

type ReactionForm struct {
	ID       int
	UserID   int
	Reaction bool
	Token    string
}

type PostForm struct {
	Title               string   `form:"title"`
	Content             string   `form:"content"`
	Categories          []int    `form:"category"`
	CategoriesString    []string `form:"category"`
	validator.Validator `form:"-"`
}

func (f *PostForm) ConverCategories(categories []string) error {
	for _, str := range f.CategoriesString {
		nb, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		if nb > len(categories) || nb < 0 {
			return UnknownCategory
		}
		f.Categories = append(f.Categories, nb)
	}
	return nil
}
