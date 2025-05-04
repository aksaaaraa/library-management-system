package repository

import (
	"database/sql"
	"errors"
	"time"

	"library-management-system/internal/models"
)

// BorrowingRepository handles database operations for borrowings
type BorrowingRepository struct {
	db *Database
}

// NewBorrowingRepository creates a new BorrowingRepository instance
func NewBorrowingRepository(db *Database) *BorrowingRepository {
	return &BorrowingRepository{db: db}
}

// GetByID retrieves a borrowing record by ID
func (r *BorrowingRepository) GetByID(id int64) (*models.Borrowing, error) {
	query := `
		SELECT b.id, b.user_id, b.book_copy_id, b.borrowed_date, b.due_date,
		       b.returned_date, b.status, b.fine_amount, b.fine_paid,
		       b.staff_id_checkout, b.staff_id_return, b.notes,
		       bc.book_id, bk.title, bk.author,
		       u.full_name as user_name
		FROM borrowings b
		JOIN book_copies bc ON b.book_copy_id = bc.id
		JOIN books bk ON bc.book_id = bk.id
		JOIN users u ON b.user_id = u.id
		WHERE b.id = ?`

	var borrowing models.Borrowing
	var returnedDate sql.NullTime
	var staffIDCheckout, staffIDReturn sql.NullInt64
	var notes sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&borrowing.ID, &borrowing.UserID, &borrowing.BookCopyID, &borrowing.BorrowedDate,
		&borrowing.DueDate, &returnedDate, &borrowing.Status, &borrowing.FineAmount,
		&borrowing.FinePaid, &staffIDCheckout, &staffIDReturn, &notes,
		&borrowing.BookID, &borrowing.BookTitle, &borrowing.BookAuthor,
		&borrowing.UserName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrBorrowingNotFound
		}
		return nil, err
	}

	if returnedDate.Valid {
		borrowing.ReturnedDate = &returnedDate.Time
	}
	if staffIDCheckout.Valid {
		borrowing.StaffIDCheckout = &staffIDCheckout.Int64
	}
	if staffIDReturn.Valid {
		borrowing.StaffIDReturn = &staffIDReturn.Int64
	}
	if notes.Valid {
		borrowing.Notes = notes.String
	}

	return &borrowing, nil
}

// Create adds a new borrowing record
func (r *BorrowingRepository) Create(tx *sql.Tx, borrowing *models.Borrowing) error {
	query := `
		INSERT INTO borrowings (
			user_id, book_copy_id, borrowed_date, due_date, status,
			staff_id_checkout, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.Exec(
			query,
			borrowing.UserID, borrowing.BookCopyID, borrowing.BorrowedDate,
			borrowing.DueDate, borrowing.Status, borrowing.StaffIDCheckout,
			borrowing.Notes,
		)
	} else {
		result, err = r.db.Exec(
			query,
			borrowing.UserID, borrowing.BookCopyID, borrowing.BorrowedDate,
			borrowing.DueDate, borrowing.Status, borrowing.StaffIDCheckout,
			borrowing.Notes,
		)
	}

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	borrowing.ID = id
	return nil
}

// UpdateBookCopyStatus updates the status of a book copy
func (r *BorrowingRepository) UpdateBookCopyStatus(tx *sql.Tx, bookCopyID int64, status string) error {
	query := `UPDATE book_copies SET status = ? WHERE id = ?`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, status, bookCopyID)
	} else {
		_, err = r.db.Exec(query, status, bookCopyID)
	}

	return err
}

// UpdateBookAvailableCopies updates the available copies count for a book
func (r *BorrowingRepository) UpdateBookAvailableCopies(tx *sql.Tx, bookID int64, increment bool) error {
	var query string
	if increment {
		query = `UPDATE books SET available_copies = available_copies + 1 WHERE id = ?`
	} else {
		query = `UPDATE books SET available_copies = available_copies - 1 WHERE id = ? AND available_copies > 0`
	}

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.Exec(query, bookID)
	} else {
		result, err = r.db.Exec(query, bookID)
	}

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 && !increment {
		return models.ErrNoAvailableCopies
	}

	return nil
}

// Return marks a borrowing as returned
func (r *BorrowingRepository) Return(borrowingID int64, returnedDate time.Time, staffID int64, fineAmount float64) error {
	// Start a transaction
	return r.db.Transaction(func(tx *sql.Tx) error {
		// Get the borrowing record
		query := `
			SELECT b.book_copy_id, bc.book_id
			FROM borrowings b
			JOIN book_copies bc ON b.book_copy_id = bc.id
			WHERE b.id = ?`

		var bookCopyID, bookID int64
		err := tx.QueryRow(query, borrowingID).Scan(&bookCopyID, &bookID)
		if err != nil {
			return err
		}

		// Update the borrowing record
		updateQuery := `
			UPDATE borrowings
			SET returned_date = ?, staff_id_return = ?, status = ?,
				fine_amount = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?`

		status := models.BorrowingStatusReturned
		if fineAmount > 0 {
			status = models.BorrowingStatusOverdue
		}

		_, err = tx.Exec(updateQuery, returnedDate, staffID, status, fineAmount, borrowingID)
		if err != nil {
			return err
		}

		// Update book copy status
		err = r.UpdateBookCopyStatus(tx, bookCopyID, models.BookCopyStatusAvailable)
		if err != nil {
			return err
		}

		// Update book available copies
		return r.UpdateBookAvailableCopies(tx, bookID, true)
	})
}

// ListByUser retrieves all borrowings for a user with pagination
func (r *BorrowingRepository) ListByUser(userID int64, page, pageSize int) ([]*models.Borrowing, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT b.id, b.book_copy_id, b.borrowed_date, b.due_date,
			b.returned_date, b.status, b.fine_amount, b.fine_paid,
			bc.book_id, bk.title, bk.author, bk.cover_image_url
		FROM borrowings b
		JOIN book_copies bc ON b.book_copy_id = bc.id
		JOIN books bk ON bc.book_id = bk.id
		WHERE b.user_id = ?
		ORDER BY 
			CASE 
				WHEN b.status = 'active' THEN 0
				WHEN b.status = 'overdue' THEN 1
				ELSE 2
			END,
			b.borrowed_date DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var borrowings []*models.Borrowing
	for rows.Next() {
		var borrowing models.Borrowing
		var returnedDate sql.NullTime

		err := rows.Scan(
			&borrowing.ID, &borrowing.BookCopyID, &borrowing.BorrowedDate,
			&borrowing.DueDate, &returnedDate, &borrowing.Status,
			&borrowing.FineAmount, &borrowing.FinePaid,
			&borrowing.BookID, &borrowing.BookTitle, &borrowing.BookAuthor,
			&borrowing.CoverImageURL,
		)
		if err != nil {
			return nil, err
		}

		borrowing.UserID = userID
		if returnedDate.Valid {
			borrowing.ReturnedDate = &returnedDate.Time
		}

		borrowings = append(borrowings, &borrowing)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return borrowings, nil
}

// CountByUser counts all borrowings for a user
func (r *BorrowingRepository) CountByUser(userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM borrowings WHERE user_id = ?`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// List retrieves all borrowings with pagination and filters
func (r *BorrowingRepository) List(page, pageSize int, filters models.BorrowingFilters) ([]*models.Borrowing, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Base query
	baseQuery := `
		SELECT b.id, b.user_id, b.book_copy_id, b.borrowed_date, b.due_date,
			b.returned_date, b.status, b.fine_amount, b.fine_paid,
			bc.book_id, bk.title, bk.author,
			u.full_name as user_name
		FROM borrowings b
		JOIN book_copies bc ON b.book_copy_id = bc.id
		JOIN books bk ON bc.book_id = bk.id
		JOIN users u ON b.user_id = u.id
		WHERE 1=1`

	// Build query with filters
	query := baseQuery
	args := []interface{}{}

	if filters.Status != "" {
		query += " AND b.status = ?"
		args = append(args, filters.Status)
	}

	if filters.BookTitle != "" {
		query += " AND bk.title LIKE ?"
		args = append(args, "%"+filters.BookTitle+"%")
	}

	if filters.UserName != "" {
		query += " AND u.full_name LIKE ?"
		args = append(args, "%"+filters.UserName+"%")
	}

	if !filters.FromDate.IsZero() {
		query += " AND b.borrowed_date >= ?"
		args = append(args, filters.FromDate)
	}

	if !filters.ToDate.IsZero() {
		query += " AND b.borrowed_date <= ?"
		args = append(args, filters.ToDate)
	}

	if filters.Overdue {
		query += " AND b.due_date < CURRENT_TIMESTAMP AND b.returned_date IS NULL"
	}

	// Add sorting
	query += " ORDER BY b.borrowed_date DESC"

	// Add pagination
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	// Execute query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var borrowings []*models.Borrowing
	for rows.Next() {
		var borrowing models.Borrowing
		var returnedDate sql.NullTime

		err := rows.Scan(
			&borrowing.ID, &borrowing.UserID, &borrowing.BookCopyID,
			&borrowing.BorrowedDate, &borrowing.DueDate, &returnedDate,
			&borrowing.Status, &borrowing.FineAmount, &borrowing.FinePaid,
			&borrowing.BookID, &borrowing.BookTitle, &borrowing.BookAuthor,
			&borrowing.UserName,
		)
		if err != nil {
			return nil, err
		}

		if returnedDate.Valid {
			borrowing.ReturnedDate = &returnedDate.Time
		}

		borrowings = append(borrowings, &borrowing)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return borrowings, nil
}

// Count returns the total number of borrowings that match the given filters
func (r *BorrowingRepository) Count(filters models.BorrowingFilters) (int, error) {
	// Base query
	baseQuery := `
		SELECT COUNT(*)
		FROM borrowings b
		JOIN book_copies bc ON b.book_copy_id = bc.id
		JOIN books bk ON bc.book_id = bk.id
		JOIN users u ON b.user_id = u.id
		WHERE 1=1`

	// Build query with filters
	query := baseQuery
	args := []interface{}{}

	if filters.Status != "" {
		query += " AND b.status = ?"
		args = append(args, filters.Status)
	}

	if filters.BookTitle != "" {
		query += " AND bk.title LIKE ?"
		args = append(args, "%"+filters.BookTitle+"%")
	}

	if filters.UserName != "" {
		query += " AND u.full_name LIKE ?"
		args = append(args, "%"+filters.UserName+"%")
	}

	if !filters.FromDate.IsZero() {
		query += " AND b.borrowed_date >= ?"
		args = append(args, filters.FromDate)
	}

	if !filters.ToDate.IsZero() {
		query += " AND b.borrowed_date <= ?"
		args = append(args, filters.ToDate)
	}

	if filters.Overdue {
		query += " AND b.due_date < CURRENT_TIMESTAMP AND b.returned_date IS NULL"
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}
