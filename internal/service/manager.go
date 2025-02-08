package service

import (
	"forum/internal/repo"
	"forum/models"
	"net/http"
)

type service struct {
	repo repo.RepoI
}

type ServiceI interface {
	UserServiceI
	CategoryServiceI
	PostServiceI
	InteractionServiceI
}

type InteractionServiceI interface {
	CommentPost(models.CommentForm) error
	PostReaction(models.ReactionForm) error
	CommentReaction(models.ReactionForm) error
	GetReactionPosts(token string) (map[int]bool, error)
	GetReactionPost(token string, postID int) (bool, bool, error)
	IsLikedPost(posts *[]models.Post, reactions map[int]bool) *[]models.Post
	IsLikedComment(posts *models.Post, reactions map[int]bool) *models.Post
	GetReactionComment(token string, postID int) (map[int]bool, error)
}

type UserServiceI interface {
	ValidToken(token string) (bool, error)
	GetUser(*http.Request) (*models.User, error)
	CreateUser(models.User) error
	Authenticate(string, string) (*models.Session, error)
	DeleteSession(string) error
}

type PostServiceI interface {
	CreatePost(string, string, string, []int) (int, error)
	GetPostByID(int) (*models.Post, error)
	GetAllPostPaginated(curentPage, pageSize int) (*[]models.Post, error)
	GetAllPostByCategoryPaginated(curentPage, pageSize, category int) (*[]models.Post, error)
	GetAllPostByCategory(category int) (*[]models.Post, error)
	GetAllPostByUserPaginated(token string, curentPage, pageSize int) (*[]models.Post, error)
	GetLikedPostsPaginated(token string, curentPage, pageSize int) (*[]models.Post, error)
	SetUpPage(data *models.TemplateData, r *http.Request) (*models.TemplateData, error)
}

type CategoryServiceI interface {
	GetAllCategory() ([]string, error)
}

func New(r repo.RepoI) ServiceI {
	return &service{
		r,
	}
}
