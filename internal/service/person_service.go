package service

import (
	"context"
	"sync"
	"uuid"

	"github.com/CakeForKit/rsoi-lab1/internal/common/custom_error"
	"github.com/CakeForKit/rsoi-lab1/internal/common/model"
	"github.com/CakeForKit/rsoi-lab1/internal/repository"
)

var personSrv PersonService
var personSrvMutex sync.Mutex

type PersonService interface {
	GetById(ctx context.Context, id uuid.UUID) (model.Person, error)
	GetAll(ctx context.Context) ([]model.Person, error)

	Create(ctx context.Context, person model.Person) (model.Person, error)
	Update(ctx context.Context, person model.Person) (model.Person, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type personService struct {
	repository repository.PersonRepository
}

func GetPersonService() (PersonService, error) {
	personSrvMutex.Lock()
	defer personSrvMutex.Unlock()

	if personSrv != nil {
		return personSrv, nil
	}

	repo, err := repository.GetPersonRepository()
	if err != nil {
		return nil, err
	}

	personSrv = &personService{repository: repo}
	return personSrv, nil
}

func (service *personService) GetById(ctx context.Context, id uuid.UUID) (model.Person, error) {
	persons, err := service.repository.GetById(ctx, []uuid.UUID{id})
	if err != nil {
		return model.Person{}, err
	}
	if len(persons) == 0 {
		return model.Person{}, custom_error.NotFoundError("person")
	}
	return persons[0], nil
}

func (service *personService) GetAll(ctx context.Context) ([]model.Person, error) {
	return service.repository.GetAll(ctx)
}

func (service *personService) Create(ctx context.Context, person model.Person) (model.Person, error) {
	persons, err := service.repository.Create(ctx, []model.Person{person})
	if err != nil {
		return model.Person{}, err
	}
	return persons[0], err
}

func (service *personService) Update(ctx context.Context, person model.Person) (model.Person, error) {
	persons, err := service.repository.Update(ctx, []model.Person{person})
	if err != nil {
		return model.Person{}, err
	}
	return persons[0], err
}

func (service *personService) DeleteById(ctx context.Context, id uuid.UUID) error {
	return service.repository.DeleteById(ctx, []uuid.UUID{id})
}
