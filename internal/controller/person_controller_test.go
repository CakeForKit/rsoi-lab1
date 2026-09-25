package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CakeForKit/rsoi-lab1/internal/common/custom_error"
	"github.com/CakeForKit/rsoi-lab1/internal/common/model"
	"github.com/CakeForKit/rsoi-lab1/internal/service"
	"github.com/gin-gonic/gin"
)

type personServiceStub struct {
	persons map[uint]model.Person
}

func newPersonServiceStub() *personServiceStub {
	return &personServiceStub{persons: make(map[uint]model.Person)}
}

func (stub *personServiceStub) GetById(_ context.Context, id uint) (model.Person, error) {
	person, ok := stub.persons[id]
	if !ok {
		return model.Person{}, custom_error.NotFoundError("person")
	}
	return person, nil
}

func (stub *personServiceStub) GetAll(_ context.Context) ([]model.Person, error) {
	persons := make([]model.Person, 0, len(stub.persons))
	for _, person := range stub.persons {
		persons = append(persons, person)
	}
	return persons, nil
}

func (stub *personServiceStub) Create(_ context.Context, person model.Person) (model.Person, error) {
	person.ID = uint(len(stub.persons) + 1)
	stub.persons[person.ID] = person
	return person, nil
}

func (stub *personServiceStub) Update(ctx context.Context, id uint, update model.PersonUpdate) (model.Person, error) {
	person, err := stub.GetById(ctx, id)
	if err != nil {
		return model.Person{}, err
	}
	if update.Name != nil {
		person.Name = *update.Name
	}
	if update.Age != nil {
		person.Age = *update.Age
	}
	if update.Address != nil {
		person.Address = *update.Address
	}
	if update.Work != nil {
		person.Work = *update.Work
	}
	stub.persons[id] = person
	return person, nil
}

func (stub *personServiceStub) DeleteById(ctx context.Context, id uint) error {
	if _, err := stub.GetById(ctx, id); err != nil {
		return err
	}
	delete(stub.persons, id)
	return nil
}

var _ service.PersonService = (*personServiceStub)(nil)

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	(&personController{service: newPersonServiceStub()}).RegisterHttpController(router)
	return router
}

func request(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}

func TestCreatePerson(t *testing.T) {
	response := request(testRouter(), http.MethodPost, "/api/v1/persons", `{"name":"Anna","age":20}`)
	if response.Code != http.StatusCreated || response.Header().Get("Location") != "/api/v1/persons/1" || response.Body.Len() != 0 {
		t.Fatalf("unexpected create response: status=%d location=%q body=%q", response.Code, response.Header().Get("Location"), response.Body.String())
	}
}

func TestGetPerson(t *testing.T) {
	router := testRouter()
	request(router, http.MethodPost, "/api/v1/persons", `{"name":"Anna","age":20}`)
	response := request(router, http.MethodGet, "/api/v1/persons/1", "")
	var person model.Person
	_ = json.Unmarshal(response.Body.Bytes(), &person)
	if response.Code != http.StatusOK || person.ID != 1 || person.Name != "Anna" {
		t.Fatalf("unexpected person response: status=%d body=%q", response.Code, response.Body.String())
	}
}

func TestGetMissingPerson(t *testing.T) {
	response := request(testRouter(), http.MethodGet, "/api/v1/persons/1", "")
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}
}

func TestPatchPerson(t *testing.T) {
	router := testRouter()
	request(router, http.MethodPost, "/api/v1/persons", `{"name":"Anna","age":20,"work":"Teacher"}`)
	response := request(router, http.MethodPatch, "/api/v1/persons/1", `{"name":"Maria"}`)
	var person model.Person
	_ = json.Unmarshal(response.Body.Bytes(), &person)
	if response.Code != http.StatusOK || person.Name != "Maria" || person.Work != "Teacher" {
		t.Fatalf("unexpected patch response: status=%d body=%q", response.Code, response.Body.String())
	}
}

func TestDeletePerson(t *testing.T) {
	router := testRouter()
	request(router, http.MethodPost, "/api/v1/persons", `{"name":"Anna"}`)
	response := request(router, http.MethodDelete, "/api/v1/persons/1", "")
	if response.Code != http.StatusNoContent || response.Body.Len() != 0 {
		t.Fatalf("unexpected delete response: status=%d body=%q", response.Code, response.Body.String())
	}
}
