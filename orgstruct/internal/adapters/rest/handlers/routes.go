package handlers

import (
	"net/http"

	"github.com/plopyblopy/orgstruct/internal/usecases"
)

// RouteRegistrar интерфейс что содержит метод для добавления обработчиков
// Он предназначен для добавления маршрута и обратчика в http.ServeMux.
type RouteRegistrar interface {
	AddHandler(pattern string, handler http.HandlerFunc)
}

func RegisterRoutes(r RouteRegistrar) {
	r.AddHandler("POST /departments", PostDepartment(usecases.PostDepartment()))
}
