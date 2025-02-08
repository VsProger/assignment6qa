package sqlite

import (
	"database/sql"
	"fmt"
)

type Sqlite struct {
	db *sql.DB
}

func NewDB(storagePath string) (*Sqlite, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tableCreationQueries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			hashed_password TEXT NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			status INTEGER DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			token TEXT NOT NULL,
			exp_time TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			like INTEGER DEFAULT 0,
			dislike INTEGER DEFAULT 0,
			image_name TEXT,
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS post_user_Like (
			user_id INTEGER,
			post_id INTEGER,
			is_like BOOLEAN,
			PRIMARY KEY (user_id, post_id),
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (post_id) REFERENCES posts(post_id)
		);`,
		`CREATE TABLE IF NOT EXISTS category (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS post_category (
			category_id INTEGER,
			post_id INTEGER,
			PRIMARY KEY (category_id, post_id),
			FOREIGN KEY (category_id) REFERENCES category(category_id),
			FOREIGN KEY (post_id) REFERENCES posts(post_id)
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY,
			post_id INTEGER,
			user_id INTEGER,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			content TEXT NOT NULL,
			like INTEGER DEFAULT 0,
			dislike INTEGER DEFAULT 0,
			FOREIGN KEY (post_id) REFERENCES posts(post_id),
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS comment_user_Like (
			user_id INTEGER,
			comment_id INTEGER,
			is_like BOOLEAN,
			PRIMARY KEY (user_id, comment_id),
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (comment_id) REFERENCES comments(comment_id)
		);`,
	}

	for _, query := range tableCreationQueries {
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		_, err = stmt.Exec()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		stmt.Close()
	}

	// defaultCategories := []string{"Technology", "Entertainment", "Sports", "Education"}
	// for _, category := range defaultCategories {
	// 	insertQuery := `INSERT INTO category (name) VALUES (?)`
	// 	stmt, err := db.Prepare(insertQuery)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("%s: %w", op, err)
	// 	}
	// 	_, err = stmt.Exec(category)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("%s: %w", op, err)
	// 	}
	// 	stmt.Close()
	// }

	return &Sqlite{db: db}, nil
}
