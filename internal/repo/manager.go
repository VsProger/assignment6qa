package repo

import (
	"forum/internal/repo/sqlite"
	"forum/models"
)

type UserRepo interface {
	CreateUser(models.User) error
	GetUserByID(int) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	UpdateUserByID(string) (*models.User, error)
	Authenticate(email, password string) (int, error)
}

type SessionRepo interface {
	GetUserIDByToken(string) (int, error)
	CreateSession(*models.Session) error
	DeleteSessionByUserID(int) error
	DeleteSessionByToken(string) error
	IsValidToken(token string) (bool, error)
}

type PostRepo interface {
	CreatePost(userID int, title, content, imageName string) (int, error)
	GetPostByID(int) (*models.Post, error)
	GetCategoriesByPostID(int) (map[int]string, error)
	// GetAllPost() (*models.Post, error)
	// UpdatePost(string, *models.Post) error
	GetLikedPostsPaginated(userID, page, pageSize int) (*[]models.Post, error)
	GetAllPostByUserIDPaginated(userID, page, pageSize int) (*[]models.Post, error)
	GetAllPostByCategory(category int) (*[]models.Post, error)
	GetPageNumber(pageSize int, category int) (int, error)
	GetAllPostPaginated(page int, pageSize int) (*[]models.Post, error)
	GetAllPostByCategoryPaginated(page int, pageSize int, category int) (*[]models.Post, error)
	GetPageNumberLikedPosts(pageSize int, userID int) (int, error)
	GetPageNumberMyPosts(pageSize int, userID int) (int, error)
	CheckPostExists(postID int) bool
}

type InteractionRepo interface {
	AddReactionPost(form models.ReactionForm) error
	DeleteReactionPost(form models.ReactionForm, isLike bool) error
	GetReactionPost(userID, postID int) (bool, bool, error)
	GetReactionPosts(userID int) (map[int]bool, error)
	GetReactionComments(userID, postID int) (map[int]bool, error)
}

type CategoryRepo interface {
	AddCategoryToPost(int, []int) error
	GetALLCategory() ([]string, error)
	// CreateCategory(string) error
}

type CommentRepo interface {
	CommentPost(models.CommentForm) error
	GetCommentsByPostID(postID int) (*[]models.Comment, error)
	// 	GetAllCommentByUserID(string) (*[]models.Post, error)
	CheckReactionComment(form models.ReactionForm) (bool, bool, error)
	AddReactionComment(form models.ReactionForm) error
	DeleteReactionComment(form models.ReactionForm, isLike bool) error
	CheckCommentExists(commentID int) bool
}

type RepoI interface {
	UserRepo
	SessionRepo
	PostRepo
	CategoryRepo
	CommentRepo
	InteractionRepo
}

func New(storagePath string) (RepoI, error) {
	return sqlite.NewDB(storagePath)
}
