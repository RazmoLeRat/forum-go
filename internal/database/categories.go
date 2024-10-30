package database

import (
	"database/sql"
	"forum-go/internal/models"
	"math"
	"strconv"

	"math/rand"
)

func (s *service) GetCategories() ([]models.Category, error) {
	rows, err := s.db.Query("SELECT * FROM Category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.CategoryId, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (s *service) AddCategory(name string) error {
	category := models.Category{
		CategoryId: sql.NullString{String: strconv.Itoa(rand.Intn(math.MaxInt32)), Valid: true},
		Name:       sql.NullString{String: name, Valid: true},
	}
	query := "INSERT INTO Category (category_id,name) VALUES (?,?)"
	_, err := s.db.Exec(query, category.CategoryId, category.Name)
	return err
}

func (s *service) DeleteCategory(id string) error {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Delete from Post_Category first to maintain referential integrity
	query := "DELETE FROM Post_Category WHERE category_id=?"
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete from Category
	query = "DELETE FROM Category WHERE category_id=?"
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *service) EditCategory(id, name string) error {
	query := "UPDATE Category SET name=? WHERE category_id=?"
	_, err := s.db.Exec(query, name, id)
	return err
}
