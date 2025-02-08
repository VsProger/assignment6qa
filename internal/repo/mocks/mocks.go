package mock

import (
	"forum/models"
	"strings"
	"testing"
)

func NewMockRepo(t *testing.T) *MockRepo {
	return &MockRepo{}
}

func Equal(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; expected: %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedStr string) {
	t.Helper()

	if !strings.Contains(actual, expectedStr) {
		t.Errorf("expected '%s' to contain '%s'", actual, expectedStr)
	}
}

type MockRepo struct{}

func (r *MockRepo) CreatePost(userID int, title, content, imageName string) (int, error) {
	return userID, nil
}

func (r *MockRepo) GetPost(id int) (*models.Post, error) {
	if id == 1 {
		return &models.Post{PostID: 1, Title: "test", Content: "test"}, nil
	}
	return nil, models.ErrNoRecord
}

func (r *MockRepo) UserPosts(userid int) ([]*models.Post, error) {
	return []*models.Post{{PostID: 1, Title: "test", Content: "test"}}, nil
}

func (r *MockRepo) LatestPosts() ([]*models.Post, error) {
	return []*models.Post{{PostID: 1, Title: "test", Content: "test"}}, nil
}

func (r *MockRepo) GetLikedPost(userid int) ([]*models.Post, error) {
	return []*models.Post{{PostID: 1, Title: "test", Content: "test"}}, nil
}

func (r *MockRepo) CreateUser(u models.User) error {
	if u.Name == "max" && u.Email == "max@gmail.com" {
		return nil
	}

	if u.Email == "max@gmail.com" {
		return models.ErrDuplicateEmail
	}
	return nil
}

func (r *MockRepo) Authenticate(email, password string) (int, error) {
	if email == "max@gmail.com" && password == "maxmax01" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (r *MockRepo) Exists(name string) (bool, error) {
	return true, nil
}

func (r *MockRepo) GetUser(id int) (string, error) {
	return "", nil
}

func (r *MockRepo) CreateReaction(userid, postid, reaction int) error {
	return nil
}

func (r *MockRepo) DeleteReactionComment(form models.ReactionForm, isLike bool) error {
	return nil
}

func (r *MockRepo) DeleteReactionPost(form models.ReactionForm, isLike bool) error {
	return nil
}

func (r *MockRepo) GetLikes(postid int) (int, error) {
	return 0, nil
}

func (r *MockRepo) GetDislikes(postid int) (int, error) {
	return 0, nil
}

func (r *MockRepo) AddReactionComment(form models.ReactionForm) error {
	return nil
}

func (r *MockRepo) AddReactionPost(form models.ReactionForm) error {
	return nil
}

func (r *MockRepo) CheckReactionComment(form models.ReactionForm) (bool, bool, error) {
	return true, true, nil
}

func (r *MockRepo) CreateComment(postid, userid int, text string) (int, error) {
	return 1, nil
}

func (r *MockRepo) GetComment(id int) (*models.Comment, error) {
	return &models.Comment{CommentID: 1, Content: "test", UserID: 1}, nil
}

func (r *MockRepo) GetComments(id int) ([]*models.Comment, error) {
	return []*models.Comment{{CommentID: 1, Content: "test", UserID: 1}}, nil
}

func (r *MockRepo) CheckCommentExists(commentID int) bool {
	return true
}

func (r *MockRepo) CheckPostExists(postID int) bool {
	return true
}

func (r *MockRepo) CommentPost(form models.CommentForm) error {
	return nil
}

func (r *MockRepo) IsValidToken(token string) (bool, error) {
	return true, nil
}

func (r *MockRepo) GetUserIDBySessionToken(sessionToken string) int {
	return 1
}

func (r *MockRepo) DeleteSessionByToken(token string) error {
	return nil
}

func (r *MockRepo) CreateSession(*models.Session) error {
	return nil
}

func (r *MockRepo) GetUserIDByToken(token string) (int, error) {
	return 1, nil
}

func (r *MockRepo) DeleteSessionByUserID(userID int) error {
	return nil
}

func (r *MockRepo) CreateCommentReaction(userid, commentid, reaction int) error {
	return nil
}

func (r *MockRepo) GetCommentLikes(commentid int) (int, error) {
	return 0, nil
}

func (r *MockRepo) GetCommentDislikes(commentid int) (int, error) {
	return 0, nil
}

func (r *MockRepo) ChooseCategories(postid int, categorie []string) error {
	return nil
}

func (r *MockRepo) AddCategoryToPost(postid int, categories []int) error {
	return nil
}

func (r *MockRepo) GetCategory(postid int) ([]string, error) {
	return []string{"1", "2", "3"}, nil
}

func (r *MockRepo) Exitsts(name string) (bool, error) {
	return true, nil
}

func (r *MockRepo) GetCategoriesByPostID(id int) (map[int]string, error) {
	return map[int]string{1: "category1", 2: "category2"}, nil
}

func (r *MockRepo) GetReactionPost(userID, postID int) (bool, bool, error) {
	if postID > 1 && postID < 1 {
		return false, false, models.ErrNoRecord
	}
	return true, true, nil
}

func (r *MockRepo) GetReactionPosts(userID int) (map[int]bool, error) {
	return map[int]bool{1: true}, nil
}

func (r *MockRepo) GetReactionComments(userID, postID int) (map[int]bool, error) {
	return map[int]bool{1: true}, nil
}

func (r *MockRepo) GetALLCategory() ([]string, error) {
	return []string{"category1", "category2"}, nil
}

func (r *MockRepo) GetPostByID(postID int) (*models.Post, error) {
	return &models.Post{
		PostID:  1,
		Title:   "test",
		Content: "test",
	}, nil
}

func (r *MockRepo) GetCommentsByPostID(postID int) (*[]models.Comment, error) {
	return &[]models.Comment{{CommentID: 1, Content: "test", UserID: 1}}, nil
}

func (s *MockRepo) GetAllPost() ([]models.Post, error) {
	return []models.Post{}, nil
}

func (s *MockRepo) GetAllPostByUserIDPaginated(userID, page, pageSize int) (*[]models.Post, error) {
	return &[]models.Post{}, nil
}

func (s *MockRepo) GetAllPostByCategory(categoryID int) (*[]models.Post, error) {
	return &[]models.Post{
		{
			PostID:    1,
			UserID:    1,
			Content:   "test",
			Title:     "test",
			Like:      0,
			Dislike:   0,
			ImageName: "test",
		},
	}, nil
}

func (s *MockRepo) GetAllPostByCategoryPaginated(page int, pageSize int, categoryID int) (*[]models.Post, error) {
	return &[]models.Post{}, nil
}

func (s *MockRepo) GetAllPostPaginated(page, pageSize int) (*[]models.Post, error) {
	return &[]models.Post{}, nil
}

func (s *MockRepo) GetLikedPostsPaginated(userID, page, pageSize int) (*[]models.Post, error) {
	return &[]models.Post{}, nil
}

func (s *MockRepo) GetPageNumber(pageSize int, category int) (int, error) {
	return 1, nil
}

func (s *MockRepo) GetPageNumberLikedPosts(pageSize int, userID int) (int, error) {
	return 1, nil
}

func (s *MockRepo) GetPageNumberMyPosts(pageSize int, userID int) (int, error) {
	return 1, nil
}

func (r *MockRepo) GetAllCommentByUserID(userID string) ([]*models.Comment, error) {
	return []*models.Comment{{CommentID: 1, Content: "test", UserID: 1}}, nil
}

func (s *MockRepo) GetUserByEmail(email string) (*models.User, error) {
	return &models.User{
		ID:    1,
		Name:  "test",
		Email: email,
	}, nil
}

func (s *MockRepo) UpdateUserByID(id string) (*models.User, error) {
	return &models.User{
		ID:    1,
		Name:  "test",
		Email: "test@example.com",
	}, nil
}

func (s *MockRepo) GetUserByID(id int) (*models.User, error) {
	return &models.User{
		ID:    1,
		Name:  "test",
		Email: "test@gmail.com",
	}, nil
}
