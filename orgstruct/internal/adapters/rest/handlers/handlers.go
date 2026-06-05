package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	. "github.com/plopyblopy/orgstruct/internal/adapters/rest"
	"github.com/plopyblopy/orgstruct/internal/domain"
)

// PostDepartment обработчик запроса на создание нового Department.
func PostDepartment(uc func(ctx context.Context, name string, parentId *int) (*domain.Department, error)) http.HandlerFunc {
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

		model, err := uc(r.Context(), i.Name, i.ParentId)
		if err != nil {
			RespondError(w, err)
			return
		}

		RespondJSON(w, http.StatusCreated, ParseDepartamentModelToResponse(*model))
	}
}

// PostEmployee обработчик запроса на создание нового Employee.
func PostEmployee(uc func(ctx context.Context, r domain.AddEmployeeRequest) (*domain.Employee, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i := struct {
			FullName string     `json:"full_name"`
			Position string     `json:"position"`
			HiredAt  *time.Time `json:"hired_at"`
		}{}

		departmentId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			RespondRowError(w, http.StatusBadRequest, "Invalid department id.")
			return
		} else if departmentId < 0 {
			RespondRowError(w, http.StatusBadRequest, "The department ID must not be less than zero.")
			return
		}

		err = json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			RespondRowError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		model, err := uc(r.Context(), domain.AddEmployeeRequest{
			FullName:     i.FullName,
			Position:     i.Position,
			DepartmentId: departmentId,
			HiredAt:      i.HiredAt,
		})
		if err != nil {
			RespondError(w, err)
			return
		}
		RespondJSON(w, http.StatusCreated, ParseEmployeeModelToResponse(*model))
	}
}

// GetDepartmentsWithChild обработчик запроса на получение Department с вложенными Employee и child Departments.
func GetDepartmentsWithChild(uc func(ctx context.Context, deptId int, depth int, includeEmployees bool) (*domain.DepartmentWithChildResponse, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		deptId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			RespondRowError(w, http.StatusBadRequest, "Invalid department id.")
			return
		} else if deptId < 0 {
			RespondRowError(w, http.StatusBadRequest, "The department ID must not be less than zero.")
			return
		}

		i := struct {
			Depth            int  `json:"depth"`
			IncludeEmployees bool `json:"include_employees"`
		}{
			Depth:            1,
			IncludeEmployees: true,
		}

		err = json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			RespondRowError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		if i.Depth < 1 || i.Depth > 5 {
			RespondRowError(w, http.StatusBadRequest, "The depth value cannot be less than 1 or more than 5.")
			return
		}

		response, err := uc(r.Context(), deptId, i.Depth, i.IncludeEmployees)
		if err != nil {
			RespondError(w, err)
			return
		}

		RespondJSON(w, http.StatusOK, response)
	}
}

// PatchDepartment обработчик запроса на обновление Department.
func PatchDepartment(uc func(ctx context.Context, dept domain.UpdateDepartment) (*domain.Department, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dept domain.UpdateDepartment

		deptId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			RespondRowError(w, http.StatusBadRequest, "Invalid department id.")
			return
		}
		dept.Id = deptId

		err = json.NewDecoder(r.Body).Decode(&dept)
		if err != nil {
			RespondRowError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		model, err := uc(r.Context(), dept)
		if err != nil {
			RespondError(w, err)
			return
		}

		RespondJSON(w, http.StatusOK, ParseDepartamentModelToResponse(*model))
	}
}
