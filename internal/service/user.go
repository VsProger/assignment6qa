package service

import (
	"forum/models"
	"forum/pkg/cookie"
	"net/http"
)

func (s *service) GetUser(r *http.Request) (*models.User, error) {
	token := cookie.GetSessionCookie(r)
	userID, err := s.repo.GetUserIDByToken(token.Value)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(userID)
}

func (s *service) DeleteSession(token string) error {
	if err := s.repo.DeleteSessionByToken(token); err != nil {
		return err
	}
	return nil
}

func (s *service) ValidToken(token string) (bool, error) {
	return s.repo.IsValidToken(token)
}

func (s *service) Authenticate(email string, password string) (*models.Session, error) {
	userID, err := s.repo.Authenticate(email, password)
	if err != nil {
		return nil, err
	}
	session := models.NewSession(userID)

	if err = s.repo.DeleteSessionByUserID(userID); err != nil {
		return nil, err
	}

	if err = s.repo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *service) CreateUser(user models.User) error {
	err := s.repo.CreateUser(user)
	return err
}
