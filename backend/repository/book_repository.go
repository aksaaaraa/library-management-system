package repository

import (
	"database/sql"
	"errors"

	"library-management-system/internal/models"
)

// BookRepository handles database operations for books
type BookRepository struct {
	db *Database
}

// NewBookRepository creates a new BookRepository instance
func NewBookRepository(db *Database) *BookRepository {
	return &BookRepository{db: db}
}

// GetByID retrieves a book by ID
func (r *BookRepository) GetByID(id int64) (*models.Book, error) {
	query := `
		SELECT id, isbn, title, author, publisher, publication_year,
		       description, category, language, page_count, total_copies,
		       available_copies, location, cover_image_url, created_at, updated_at
		FROM books
		WHERE id = ?`

	var book models.Book
	err := r.db.QueryRow(query, id).Scan(
		&book.ID, &book.ISBN, &book.Title, &book.Author, &book.Publisher,
		&book.PublicationYear, &book.Description, &book.Category, &book.Language,
		&book.PageCount, &book.TotalCopies, &book.AvailableCopies, &book.Location,
		&book.CoverImageURL, &book.CreatedAt, &book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrBookNotFound
		}
		return nil, err
	}

	return &book, nil
}

// GetByISBN retrieves a book by ISBN
func (r *BookRepository) GetByISBN(isbn string) (*models.Book, error) {
	query := `
		SELECT id, isbn, title, author, publisher, publication_year,
			   description, category, language, page_count, total_copies,
			   available_copies, location, cover_image_url, created_at, updated_at
		FROM books
		WHERE isbn = ?`

	var book models.Book
	err := r.db.QueryRow(query, isbn).Scan(
		&book.ID, &book.ISBN, &book.Title, &book.Author, &book.Publisher,
		&book.PublicationYear, &book.Description, &book.Category, &book.Language,
		&book.PageCount, &book.TotalCopies, &book.AvailableCopies, &book.Location,
		&book.CoverImageURL, &book.CreatedAt, &book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrBookNotFound
		}
		return nil, err
	}

	return &book, nil
}
