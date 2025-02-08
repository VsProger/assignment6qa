package sqlite

import "fmt"

func (s *Sqlite) AddCategoryToPost(postID int, categories []int) error {
	const op = "sqlite.AddCategoryToPost"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare("INSERT INTO post_category (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, categoryID := range categories {
		_, err = stmt.Exec(postID, categoryID)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("%s: exec statement: %w", op, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (s *Sqlite) GetALLCategory() ([]string, error) {
	op := "sqlite.GetAllCategory"
	stmt := `SELECT name FROM category ORDER BY id ASC`

	rows, err := s.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var categoryName string
		err := rows.Scan(&categoryName)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		categories = append(categories, categoryName)
	}

	return categories, nil
}

func CreateCategory(string) error {
	return nil
}

func (s *Sqlite) GetCategoriesByPostID(postID int) (map[int]string, error) {
	stmt := `SELECT 
	category_id, 
	category.name as name
	FROM 
	post_category 
	INNER JOIN category ON post_category.category_id = category.id
	WHERE post_id=?`

	rows, err := s.db.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	category := make(map[int]string)
	for rows.Next() {
		var categoryID int
		var categoryName string
		err := rows.Scan(&categoryID, &categoryName)
		if err != nil {
			return nil, err
		}
		category[categoryID] = categoryName
	}
	return category, nil
}
