package handlers

import (
	"encoding/json"
	"net/http"

	. "github.com/plopyblopy/orgstruct/internal/adapters/rest"
	"github.com/plopyblopy/orgstruct/internal/domain"
)

// PostDepartment обработчик запроса на создание нового Department
func PostDepartment(uc func(name string, parentId *int) (*domain.Department, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i := struct {
			Name     string `json:"name"`
			ParentId *int   `json:"parent_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			RespondRowError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		model, err := uc(i.Name, i.ParentId)
		if err != nil {
			RespondError(w, err)
			return
		}

		RespondJSON(w, http.StatusCreated, model)
	}
}
