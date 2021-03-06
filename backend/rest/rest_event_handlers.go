package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"net/http"
	"time"
)

type eventHandlers struct {
	validator *FormValidator
	service   EventService
}

type EventService interface {
	Create(record repository.Event, userID string) (string, error)
	Update(event repository.Event, userID string) (string, error)
	Delete(eventID string, userID string) error
	GetUsedNames(userID string, tagID string, from time.Time, to time.Time) ([]string, error)
}

func (t *eventHandlers) create(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	var event repository.Event
	err = json.NewDecoder(request.Body).Decode(&event)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = t.validator.validate(&event)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	event.PrepareReceived()

	id, err := t.service.Create(event, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't create event", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *eventHandlers) update(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	id := chi.URLParam(request, "id")
	if len(id) == 0 {
		SendErrorJSON(writer, request, http.StatusBadRequest, errors.New("id can't be empty"), "", ErrInternal)
		return
	}

	var event repository.Event
	err = json.NewDecoder(request.Body).Decode(&event)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = t.validator.validate(&event)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	event.PrepareReceived()
	event.ID = id

	id, err = t.service.Update(event, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update event", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *eventHandlers) delete(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	id := chi.URLParam(request, "id")

	err = t.service.Delete(id, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete event", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *eventHandlers) getUsedNames(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	tagID := chi.URLParam(request, "tagID")
	if len(tagID) == 0 {
		SendErrorJSON(writer, request, http.StatusBadRequest, errors.New("tag id can't be empty"), "", ErrInternal)
		return
	}

	now := time.Now()

	names, err := t.service.GetUsedNames(user.ID, tagID, now.AddDate(0, 0, -14), now)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't get events", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, names)
}
