package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *appDependencies) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFound)
	router.MethodNotAllowed = http.HandlerFunc(a.notAllowed)
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", a.healthCheckHandler)

	router.HandlerFunc(http.MethodGet, "/api/v1/books", a.requireActivatedUser(a.listBooksHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:id", a.requireActivatedUser(a.displayBookHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/search", a.requireActivatedUser(a.searchBooksHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/lists", a.requireActivatedUser(a.listListsHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/lists/:id", a.requireActivatedUser(a.displayListHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/books/:id/reviews", a.requireActivatedUser(a.displayReviewHandler))
	// router.HandlerFunc(http.MethodGet, "/api/v1/users/:id", a.requireActivatedUser(a.displayUserHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:id/lists", a.requireActivatedUser(a.displayUserListsHandler))
	// router.HandlerFunc(http.MethodGet, "/api/v1/users/:id/reviews", a.requireActivatedUser(a.displayUserReviewsHandler))

	router.HandlerFunc(http.MethodPost, "/api/v1/users", a.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/books", a.requireActivatedUser(a.createBookHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/lists", a.requireActivatedUser(a.createListHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/lists/:id/books", a.requireActivatedUser(a.addBookToListHandler))
	router.HandlerFunc(http.MethodPost, "/api/vi/books/:id/reviews", a.requireActivatedUser(a.createReviewHandler))

	router.HandlerFunc(http.MethodPut, "/api/v1/users/activated", a.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/books/:id", a.requireActivatedUser(a.updateBookHandler))
	router.HandlerFunc(http.MethodPut, "/api/v1/lists/:id", a.requireActivatedUser(a.updateListHandler))
	router.HandlerFunc(http.MethodPut, "/api/v1/reviews/:id", a.requireActivatedUser(a.updateReviewHandler))

	router.HandlerFunc(http.MethodDelete, "/api/v1/books/:id", a.requireActivatedUser(a.deleteBookHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:id", a.requireActivatedUser(a.deleteListHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/lists/:id/books", a.requireActivatedUser(a.deleteBookFromListHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/reviews/:id", a.requireActivatedUser(a.deleteReviewHandler))

	return a.recoverPanic(a.rateLimit(a.authenticate(router)))
}
