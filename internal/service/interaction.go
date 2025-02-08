package service

import (
	"forum/models"
)

func (s *service) CommentPost(form models.CommentForm) error {
	var err error
	form.UserID, err = s.repo.GetUserIDByToken(form.Token)
	if err != nil {
		return err
	}
	return s.repo.CommentPost(form)
}

func (s *service) PostReaction(form models.ReactionForm) error {
	var err error
	form.UserID, err = s.repo.GetUserIDByToken(form.Token)
	if err != nil {
		return err
	}
	ok := s.repo.CheckPostExists(form.ID)
	if !ok {
		return models.ErrNoRecord
	}
	exists, isLike, err := s.repo.GetReactionPost(form.UserID, form.ID)
	if err != nil {
		return err
	}
	if exists {
		err := s.repo.DeleteReactionPost(form, isLike)
		if err != nil {
			return err
		}
		if isLike == form.Reaction {
			return nil
		}
	}

	err = s.repo.AddReactionPost(form)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CommentReaction(form models.ReactionForm) error {
	var err error
	form.UserID, err = s.repo.GetUserIDByToken(form.Token)
	if err != nil {
		return err
	}
	ok := s.repo.CheckCommentExists(form.ID)

	if !ok {
		return models.ErrNoRecord
	}

	exists, isLike, err := s.repo.CheckReactionComment(form)
	if err != nil {
		return err
	}
	if exists {
		err := s.repo.DeleteReactionComment(form, isLike)
		if err != nil {
			return err
		}
		if isLike == form.Reaction {
			return nil
		}
	}

	err = s.repo.AddReactionComment(form)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetReactionPosts(token string) (map[int]bool, error) {
	userID, err := s.repo.GetUserIDByToken(token)
	if err != nil {
		return nil, err
	}
	reactions, err := s.repo.GetReactionPosts(userID)
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (s *service) GetReactionPost(token string, postID int) (bool, bool, error) {
	userID, err := s.repo.GetUserIDByToken(token)
	if err != nil {
		return false, false, err
	}

	exists, reaction, err := s.repo.GetReactionPost(userID, postID)
	if err != nil {
		return false, false, err
	}

	return exists, reaction, nil
}

func (s *service) GetReactionComment(token string, postID int) (map[int]bool, error) {
	userID, err := s.repo.GetUserIDByToken(token)
	if err != nil {
		return nil, err
	}

	reactions, err := s.repo.GetReactionComments(userID, postID)
	if err != nil {
		return nil, err
	}

	return reactions, nil
}

func (s *service) IsLikedPost(posts *[]models.Post, reactions map[int]bool) *[]models.Post {
	postCopy := *posts
	for key, value := range reactions {
		for i, post := range postCopy {
			if post.PostID == key {
				if value == true {
					postCopy[i].IsLiked = 1
				} else {
					postCopy[i].IsLiked = -1
				}

				break
			}
		}
	}
	return &postCopy
}

func (s *service) IsLikedComment(post *models.Post, reactions map[int]bool) *models.Post {
	postCopy := *post
	if post.Comment != nil {
		newComments := *post.Comment

		if len(newComments) > 0 && newComments[0].Content != "" {
			for key, value := range reactions {
				for i, comment := range *postCopy.Comment {
					if comment.CommentID == key {
						if value == true {
							newComments[i].IsLiked = 1
						} else {
							newComments[i].IsLiked = -1
						}
						break
					}
				}
			}
		}

		postCopy.Comment = &newComments
	}

	return &postCopy
}
