package controller

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"time"
	"userapiREF/cases"
	"userapiREF/entity"
)

type Controller struct {
	usecase cases.Usecase
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

var (
	UserNotFound = errors.New("user_not_found")
)

func NewController(usecase cases.Usecase) *Controller {
	return &Controller{
		usecase: usecase,
	}
}

func Build(r *chi.Mux, usecase cases.Usecase) {
	ctr := NewController(usecase)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", ctr.searchUsers)
				r.Post("/", ctr.createUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", ctr.getUser)
					r.Patch("/", ctr.updateUser)
					r.Delete("/", ctr.deleteUser)
				})
			})
		})
	})
}

func (s *Controller) searchUsers(w http.ResponseWriter, r *http.Request) {
	data, err := s.usecase.SearchUsers()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("not enough users"))
	}
	resp, _ := json.MarshalIndent(data.List, "", "  ")
	render.Status(r, 200)
	w.Write(resp)
}

func (s *Controller) createUser(w http.ResponseWriter, r *http.Request) {
	u := &entity.User{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		log.Fatal(err)
	}
	i := s.usecase.CreateUser(u)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": i,
	})
}

func (s *Controller) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data := s.usecase.GetUser(id)
	if data == nil {
		w.WriteHeader(400)
		w.Write([]byte("User with this id not found"))
		return
	}
	resp, _ := json.MarshalIndent(data, "", "  ")
	w.Write(resp)
}

func (s *Controller) updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	u := &entity.User{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		log.Fatal(err)
	}
	err1 := s.usecase.UpdateUser(u, id)
	if err1 != nil {
		_ = render.Render(w, r, ErrInvalidRequest(UserNotFound))
		return
	}
	render.Status(r, 204)
	w.Write([]byte("action done"))
}

func (s *Controller) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := s.usecase.DeleteUser(id)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(UserNotFound))
		return
	}
	render.Status(r, 204)
	w.Write([]byte("user was deleted"))
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
