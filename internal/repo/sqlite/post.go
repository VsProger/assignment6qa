package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
)

func (s *Sqlite) CheckPostExists(postID int) bool {
	var isExists bool
	checkQuery := `SELECT EXISTS(SELECT id FROM posts WHERE id = ?)`
	err := s.db.QueryRow(checkQuery, postID).Scan(&isExists)
	if err != nil {
		return false
	}
	return isExists
}

func (s *Sqlite) CreatePost(userID int, title, content, imageName string) (int, error) {
	op := "sqlite.CreatePost"
	const query = `INSERT INTO posts (user_id, title, content, image_name) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, userID, title, content, imageName)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return int(postID), nil
}

func (s *Sqlite) GetPostByID(postID int) (*models.Post, error) {
	op := "sqlite.GetPostByID"
	stmt := `SELECT p.id, p.user_id, p.title, p.content, p.created, p.like, p.dislike, p.image_name, u.name
	FROM posts p
	JOIN users u ON p.user_id = u.id 
	WHERE p.id = ?
`
	post := models.Post{}

	err := s.db.QueryRow(stmt, postID).Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName, &post.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &post, nil
}

func (s *Sqlite) GetAllPost() ([]models.Post, error) {
	const query = `SELECT post_id, user_id, title, content, created, like, dislike, image_name FROM Post`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Sqlite) GetAllPostByUserIDPaginated(userID, page, pageSize int) (*[]models.Post, error) {
	offset := (page - 1) * pageSize
	const query = `SELECT p.id, p.user_id, p.title, p.content, p.created, p.like, p.dislike, p.image_name, u.name, (SELECT COUNT(*) FROM comments c WHERE c.post_id=p.id) 
	FROM posts p 
	JOIN users u ON p.user_id = u.id
	WHERE p.user_id = ?
	ORDER BY p.created DESC
	LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName, &post.UserName, &post.CommentCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (s *Sqlite) GetAllPostByCategory(categoryID int) (*[]models.Post, error) {
	query := `SELECT p.id, p.user_id, p.title, p.content, p.created, p.like, p.dislike, p.image_name
              FROM posts AS p
              INNER JOIN post_category AS pc ON p.id = pc.post_id
              WHERE pc.category_id IN (?)
              GROUP BY p.id`

	rows, err := s.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (s *Sqlite) GetAllPostByCategoryPaginated(page int, pageSize int, categoryID int) (*[]models.Post, error) {
	// op := "sqlite.GetAllPostByCategoryPaginated"
	offset := (page - 1) * pageSize
	query := `SELECT p.id, p.user_id, p.title, p.content, p.created, p.like, p.dislike, p.image_name, u.name, (SELECT COUNT(*) FROM comments c WHERE c.post_id=p.id)
              FROM posts AS p
              INNER JOIN post_category AS pc ON p.id = pc.post_id
			  JOIN users u ON p.user_id = u.id 
              WHERE pc.category_id IN (?)
              GROUP BY p.id
			  ORDER BY p.created DESC
			  LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, categoryID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName, &post.UserName, &post.CommentCount); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (s *Sqlite) GetAllPostPaginated(page, pageSize int) (*[]models.Post, error) {
	op := "sqlite.GetAllPostPaginated"
	offset := (page - 1) * pageSize
	// stmt := `SELECT p.id, p.user_id, p.title, p.content, p.created, p.like, p.dislike, p.image_name, u.name, COUNT(c.id)
	// FROM posts p
	// JOIN users u ON p.user_id = u.id
	// right JOIN comments c on p.id = c.post_id
	// ORDER BY p.created DESC
	// LIMIT ? OFFSET ?
	// `

	stmt := `SELECT p.id, p.user_id, p.title, p.content, p.created, p.like, p.dislike, p.image_name, u.name, (SELECT COUNT(*) FROM comments c WHERE c.post_id=p.id)
	FROM posts p 
	Inner JOIN users u ON p.user_id = u.id 
	ORDER BY p.created DESC
	LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(stmt, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName, &post.UserName, &post.CommentCount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func (s *Sqlite) GetLikedPostsPaginated(userID, page, pageSize int) (*[]models.Post, error) {
	offset := (page - 1) * pageSize
	const query = `SELECT p.id, p.user_id, p.title, p.content, p.created, p.like, p.dislike, p.image_name, u.name, (SELECT COUNT(*) FROM comments c WHERE c.post_id=p.id) 
	FROM posts p 
	JOIN users u ON p.user_id = u.id
	JOIN post_user_Like l ON p.id = l.post_id
	WHERE l.user_id = ? AND l.is_like = TRUE
	GROUP BY p.id
	ORDER BY p.created DESC
	LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName, &post.UserName, &post.CommentCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (s *Sqlite) GetPageNumber(pageSize int, category int) (int, error) {
	var totalPosts int
	op := "sqlite.GetPageNumber"
	if category == 0 {
		stmt := `SELECT COUNT(*) FROM posts`
		err := s.db.QueryRow(stmt).Scan(&totalPosts)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		stmt := `SELECT COUNT (*)
			FROM posts AS p
			INNER JOIN post_category AS pc ON p.id = pc.post_id
			WHERE pc.category_id = (?)
			`
		err := s.db.QueryRow(stmt, category).Scan(&totalPosts)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}

	}

	totalPages := (totalPosts + pageSize - 1) / pageSize
	return totalPages, nil
}

func (s *Sqlite) GetPageNumberLikedPosts(pageSize int, userID int) (int, error) {
	var totalPosts int
	op := "sqlite.GetPageNumberLikedPosts"

	stmt := `SELECT COUNT(*)
	FROM posts p 
	JOIN users u ON p.user_id = u.id
	JOIN post_user_Like l ON p.id = l.post_id
	WHERE l.user_id = ? AND l.is_like = TRUE
	`
	err := s.db.QueryRow(stmt, userID).Scan(&totalPosts)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	totalPages := (totalPosts + pageSize - 1) / pageSize
	return totalPages, nil
}

func (s *Sqlite) GetPageNumberMyPosts(pageSize int, userID int) (int, error) {
	var totalPosts int
	op := "sqlite.GetPageNumberMyPosts"

	stmt := `SELECT COUNT(*) 
	FROM posts p 
	JOIN users u ON p.user_id = u.id
	WHERE p.user_id = ?
	`
	err := s.db.QueryRow(stmt, userID).Scan(&totalPosts)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	totalPages := (totalPosts + pageSize - 1) / pageSize
	return totalPages, nil
}
