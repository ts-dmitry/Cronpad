package rest

import (
	"github.com/go-chi/render"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"github.com/ts-dmitry/cronpad/backend/service"
	"net/http"
)

type projectHandlers struct {
	store       ProjectStore
	userService *service.UserService
}

type ProjectStore interface {
	FindAllActiveProjectsByUser(userID string) (repository.Projects, error)
}

func (t *projectHandlers) findAllByUser(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	projects, err := t.store.FindAllActiveProjectsByUser(user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get project", ErrInternal)
		return
	}

	for i := range projects {
		projects[i].PrepareToSend()
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, projects)
}
