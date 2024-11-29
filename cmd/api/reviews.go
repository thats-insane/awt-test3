package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/thats-insane/awt-test3/internal/data"
	"github.com/thats-insane/awt-test3/internal/validator"
)

func (a *appDependencies) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		BookID    int64     `json:"book_id"`
		UserID    int64     `json:"user_id"`
		Rating    int64     `json:"rating"`
		Desc      string    `json:"desc"`
		CreatedAt time.Time `json:"-"`
	}
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequest(w, r, err)
		return
	}

	review := &data.Review{
		BookID:    incomingData.BookID,
		UserID:    incomingData.UserID,
		Rating:    incomingData.Rating,
		Desc:      incomingData.Desc,
		CreatedAt: incomingData.CreatedAt,
	}
	v := validator.New()
	data.ValidateReview(v, review)
	if !v.IsEmpty() {
		a.failedValidation(w, r, v.Errors)
		return
	}

	err = a.reviewModel.Insert(review)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("api/v1/review/%d", review.ID))
	data := envelope{
		"review": review,
	}

	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)
}

func (a *appDependencies) displayReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFound(w, r)
		return
	}

	review, err := a.reviewModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFound(w, r)
		default:
			a.serverErr(w, r, err)
		}
		return
	}

	data := envelope{
		"review": review,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}
}

func (a *appDependencies) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFound(w, r)
		return
	}

	review, err := a.reviewModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFound(w, r)
		default:
			a.serverErr(w, r, err)
		}
		return
	}

	var incomingData struct {
		BookID    *int64     `json:"book_id"`
		UserID    *int64     `json:"user_id"`
		Rating    *int64     `json:"rating"`
		Desc      *string    `json:"desc"`
		CreatedAt *time.Time `json:"-"`
	}
	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequest(w, r, err)
		return
	}

	if incomingData.BookID != nil {
		review.BookID = *incomingData.BookID
	}
	if incomingData.UserID != nil {
		review.UserID = *incomingData.UserID
	}
	if incomingData.Rating != nil {
		review.Rating = *incomingData.Rating
	}
	if incomingData.Desc != nil {
		review.Desc = *incomingData.Desc
	}
	if incomingData.CreatedAt != nil {
		review.CreatedAt = *incomingData.CreatedAt
	}

	v := validator.New()
	data.ValidateReview(v, review)
	if !v.IsEmpty() {
		a.failedValidation(w, r, v.Errors)
		return
	}

	err = a.reviewModel.Update(review)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}

	data := envelope{
		"review": review,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}
}

func (a *appDependencies) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFound(w, r)
		return
	}

	err = a.reviewModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFound(w, r)
		default:
			a.serverErr(w, r, err)
		}
		return
	}

	data := envelope{
		"message": "review successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErr(w, r, err)
	}
}

/*
func (a *appDependencies) listReviewsHandler(w http.ResponseWriter, r *http.Request) {
	var queryParametersData struct {
		// Product string
		data.Filters
	}

	queryParameters := r.URL.Query()
	// queryParametersData.Product = a.getSingleQueryParameters(queryParameters, "product", "")

	queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "id")
	// queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "rating")
	// queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "helpful_count")
	// queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "created_at")
	// queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "updated_at")

	// queryParametersData.Filters.SortSafeList = []string{"id", "rating", "helpful_count", "created_at", "updated_at", "-id", "-rating", "-helpful_count", "-created_at", "-updated_at"}
	queryParametersData.Filters.SortSafeList = []string{"id", "-id"}

	v := validator.New()

	queryParametersData.Filters.Page = a.getSingleIntegerParameters(queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameters(queryParameters, "page_size", 10, v)

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidation(w, r, v.Errors)
		return
	}

	// product_id, err := toInt(queryParametersData.Product)

	// if err != nil {
	// 	a.serverErr(w, r, err)
	// 	return
	// }

	review, err := a.bookclub.GetAllReviews(queryParametersData.Filters)

	if err != nil {
		a.serverErr(w, r, err)
		return
	}

	data := envelope{
		"review": review,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		a.serverErr(w, r, err)
	}
}
*/
