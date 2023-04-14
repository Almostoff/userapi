package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"userapiREF/cases"
	"userapiREF/controller"
	"userapiREF/repository"
)

func Run() {
	repository := repository.NewRepository()
	usecase := cases.NewUseCase(repository)
	r := chi.NewRouter()
	controller.Build(r, usecase)
	http.ListenAndServe(":3333", r)
}
