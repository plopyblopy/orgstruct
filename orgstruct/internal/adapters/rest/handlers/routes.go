package handlers

import (
	"net/http"

	"github.com/plopyblopy/orgstruct/internal/adapters/postgres"
	"github.com/plopyblopy/orgstruct/internal/usecases"
)

// RouteRegistrar интерфейс что содержит метод для добавления обработчиков
// Он предназначен для добавления маршрута и обратчика в http.ServeMux.
type RouteRegistrar interface {
	AddHandler(pattern string, handler http.HandlerFunc)
}

// RegisterRoutes регистрирует маршруты, создает зависимости.
func RegisterRoutes(r RouteRegistrar, db *postgres.Db) {
	departamentRepo := postgres.NewDepartmentRepository(db)
	employeeRepo := postgres.NewEmployeeRepository(db)

	r.AddHandler("POST /departments", PostDepartment(usecases.PostDepartment(departamentRepo)))
	r.AddHandler("POST /departments/{id}/employees", PostEmployee(usecases.PostEmployee(employeeRepo)))
}
