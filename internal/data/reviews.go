package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/thats-insane/awt-test3/internal/validator"
)

type Review struct {
	ID        int64     `json:"id"`
	BookID    int64     `json:"book_id"`
	UserID    int64     `json:"user_id"`
	Rating    int64     `json:"rating"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"-"`
}

type ReviewModel struct {
	DB *sql.DB
}

func (r ReviewModel) Insert(review *Review) error {
	query := `
	INSERT INTO reviews (book_id, user_id, rating, desc)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
	`
	args := []any{review.BookID, review.UserID, review.Rating, review.Desc}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID)
}

func (r ReviewModel) Get(id int64) (*Review, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, book_id, user_id, rating, desc
	FROM reviews
	WHERE id = $1
	`
	var review Review
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(&review.ID, &review.BookID, &review.UserID, &review.Rating, &review.Desc, &review.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &review, nil
}

func (r ReviewModel) GetAll(filters Filters) ([]*Review, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT id, book_id, user_id, rating, desc
	FROM reviews
	ORDER BY %s %s, id ASC
	LIMIT $1 OFFSET $2
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var totalRecords int
	reviews := []*Review{}

	for rows.Next() {
		var review Review
		err := rows.Scan(&totalRecords, &review.ID, &review.BookID, &review.UserID, &review.Rating, &review.Desc, &review.CreatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}
		reviews = append(reviews, &review)
	}

	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)

	return reviews, metadata, nil
}

func (r ReviewModel) Update(review *Review) error {
	query := `
		UPDATE reviews
		SET book_id = $1, user_id = $2, rating = $3
		WHERE id = $4
		RETURNING id
	`

	args := []any{review.BookID, review.UserID, review.Rating, review.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID)
}

func (r ReviewModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM reviews
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateReview(v *validator.Validator, review *Review) {
	v.Check(review.BookID > 0, "review", "must be a positive integer")
	v.Check(review.UserID > 0, "review", "must be a positive integer")
	v.Check(review.Desc != "", "review", "must be provided")
	v.Check(len(review.Desc) <= 225, "review", "must not be more than 225 bytes long")
	v.Check(review.Rating >= 1 && review.Rating <= 5, "review", "must be between 1 and 5")
}