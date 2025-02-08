package sqlite

import (
	"fmt"
	"forum/models"
)

func (s *Sqlite) GetReactionPost(userID, postID int) (bool, bool, error) {
	// Check if the user has already liked/disliked the post
	var isExists bool
	checkQuery := `SELECT EXISTS(SELECT is_like FROM Post_User_Like WHERE user_id = ? AND post_id = ?)`
	err := s.db.QueryRow(checkQuery, userID, postID).Scan(&isExists)
	if err != nil {
		return false, false, err
	}
	var dbLike bool
	if isExists {
		checkQuery = `SELECT is_like FROM Post_User_Like WHERE user_id = ? AND post_id = ?`
		err = s.db.QueryRow(checkQuery, userID, postID).Scan(&dbLike)
		if err != nil {
			return false, false, err
		}
	}

	return isExists, dbLike, nil
}

func (s *Sqlite) AddReactionPost(form models.ReactionForm) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Insert like/dislike
	insertQuery := `INSERT INTO Post_User_Like (user_id, post_id, is_like) VALUES (?, ?, ?)`
	_, err = tx.Exec(insertQuery, form.UserID, form.ID, form.Reaction)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Post like/dislike count
	updateQuery := ""
	if form.Reaction {
		updateQuery = `UPDATE Posts SET like = like + 1 WHERE id = ?`
	} else {
		updateQuery = `UPDATE Posts SET dislike = dislike + 1 WHERE id = ?`
	}
	_, err = tx.Exec(updateQuery, form.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Sqlite) DeleteReactionPost(form models.ReactionForm, isLike bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// delete the like/dislike
	deleteQuery := `DELETE FROM Post_User_Like WHERE user_id = ? AND post_id = ?`
	_, err = tx.Exec(deleteQuery, form.UserID, form.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// decrement the like or dislike
	updateQuery := ""
	if isLike {
		updateQuery = `UPDATE Posts SET like = like - 1 WHERE id = ? AND like > 0`
	} else {
		updateQuery = `UPDATE Posts SET dislike = dislike - 1  WHERE id = ? AND dislike > 0`
	}
	_, err = tx.Exec(updateQuery, form.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Sqlite) GetReactionComments(userID int, postID int) (map[int]bool, error) {
	op := "sqlite.GetReactionPost"

	stmt := `SELECT comment_id, is_like FROM Comment_User_Like WHERE user_id = ?`

	rows, err := s.db.Query(stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	reactions := make(map[int]bool)
	for rows.Next() {
		var comment int
		var react bool
		if err := rows.Scan(&comment, &react); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reactions[comment] = react
	}

	return reactions, nil
}

func (s *Sqlite) GetReactionPosts(userID int) (map[int]bool, error) {
	op := "sqlite.GetReactionPosts"

	stmt := `SELECT post_id, is_like FROM Post_User_Like WHERE user_id = ?`

	rows, err := s.db.Query(stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	reactions := make(map[int]bool)
	for rows.Next() {
		var post int
		var react bool
		if err := rows.Scan(&post, &react); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reactions[post] = react
	}
	return reactions, nil
}


func (s *Sqlite) CheckReactionComment(form models.ReactionForm) (bool, bool, error) {
	// Check if the user has already liked/disliked the post
	var isExists bool
	checkQuery := `SELECT EXISTS(SELECT is_like FROM Comment_User_Like WHERE user_id = ? AND comment_id = ?)`
	err := s.db.QueryRow(checkQuery, form.UserID, form.ID).Scan(&isExists)
	if err != nil {
		return false, false, err
	}
	var dbLike bool
	if isExists {
		checkQuery = `SELECT is_like FROM Comment_User_Like WHERE user_id = ? AND comment_id = ?`
		err = s.db.QueryRow(checkQuery, form.UserID, form.ID).Scan(&dbLike)
		if err != nil {
			return false, false, err
		}
	}

	return isExists, dbLike, nil
}