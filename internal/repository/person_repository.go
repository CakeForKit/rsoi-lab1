package repository

import (
	"sync"

	"github.com/CakeForKit/rsoi-lab1/internal/common/model"
	coreRepository "github.com/CakeForKit/rsoi-lab1/internal/common/repository"
)

var personRepo *personRepository
var personRepoMutex sync.Mutex

type PersonRepository interface {
	coreRepository.Repository[model.Person]
}

type personRepository struct {
	coreRepository.Repository[model.Person]
}

func GetPersonRepository() (PersonRepository, error) {
	personRepoMutex.Lock()
	defer personRepoMutex.Unlock()

	if personRepo != nil {
		return personRepo, nil
	}

	repository, err := coreRepository.NewPostgresRepository[model.Person]()
	if err != nil {
		return nil, err
	}
	personRepo = &personRepository{repository}

	return personRepo, nil
}
