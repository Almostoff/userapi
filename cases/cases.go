package cases

import (
	"userapiREF/entity"
)

type (
	Usecase interface {
		SearchUsers() (*entity.UserStore, error)
		CreateUser(data *entity.User) string
		GetUser(id string) *entity.User
		UpdateUser(data *entity.User, id string) error
		DeleteUser(id string) error
	}

	Repository interface {
		SearchUsers() (*entity.UserStore, error)
		CreateUser(data *entity.User) string
		GetUser() *entity.UserStore
		UpdateUser(data *entity.User, id string) error
		DeleteUser(id string) error
	}
)

type usecase struct {
	repository Repository
}

func NewUseCase(repository Repository) *usecase {
	return &usecase{
		repository: repository,
	}
}

func (u *usecase) SearchUsers() (*entity.UserStore, error) {
	data, err := u.repository.SearchUsers()
	return data, err
}

func (u *usecase) CreateUser(data *entity.User) string {
	res := u.repository.CreateUser(data)
	return res
}

func (u *usecase) GetUser(id string) *entity.User {
	data := u.repository.GetUser()
	return data.List[id]
}

func (u *usecase) UpdateUser(data *entity.User, id string) error {
	err := u.repository.UpdateUser(data, id)
	return err
}

func (u *usecase) DeleteUser(id string) error {
	err := u.repository.DeleteUser(id)
	return err
}
