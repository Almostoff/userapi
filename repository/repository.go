package repository

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"strconv"
	"time"
	"userapiREF/entity"
)

const rep = `repository/users.json`

type repository struct {
	users entity.UserStore
}

func NewRepository() *repository {
	return &repository{
		users: entity.UserStore{0, make(map[string]*entity.User)},
	}
}

func (r repository) SearchUsers() (*entity.UserStore, error) {
	f, _ := ioutil.ReadFile(rep)
	s := &r.users
	err := json.Unmarshal(f, &s)
	if err != nil {
		log.Fatal("Empty file. No one users was added... ", err)
		return nil, err
	}
	return s, nil
}

func (r repository) CreateUser(data *entity.User) string {
	f, _ := ioutil.ReadFile(rep)
	s := &r.users
	_ = json.Unmarshal(f, &s)
	s.Increment++
	u := entity.User{
		CreatedAt:   time.Now(),
		DisplayName: data.DisplayName,
		Email:       data.Email,
	}

	id := strconv.Itoa(s.Increment)
	s.List[id] = &u

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(rep, b, fs.ModePerm)
	return id
}

func (r repository) GetUser() *entity.UserStore {
	f, _ := ioutil.ReadFile(rep)
	s := &r.users
	_ = json.Unmarshal(f, &s)
	return s
}

func (r repository) UpdateUser(data *entity.User, id string) error {
	f, err := ioutil.ReadFile(rep)
	s := &r.users
	err = json.Unmarshal(f, &s)

	if _, ok := s.List[id]; !ok {
		return err
	}

	u := s.List[id]
	u.DisplayName = data.DisplayName
	u.Email = data.Email
	s.List[id] = u

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(rep, b, fs.ModePerm)
	return nil
}

func (r repository) DeleteUser(id string) error {
	f, err := ioutil.ReadFile(rep)
	s := &r.users
	err = json.Unmarshal(f, &s)

	if _, ok := s.List[id]; !ok {
		return err
	}

	delete(s.List, id)

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(rep, b, fs.ModePerm)

	return nil
}
