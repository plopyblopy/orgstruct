package rest

import (
	"fmt"
	"net/http"
)

// Router оборачивает стандартный *http.Router, предоставляя методы
// для группировки маршрутов и управления префиксами.
type Router struct {
	mux *http.ServeMux
}

// NewRouter создаёт новый ServeMux с зарегистрированными обработчиками /live, /ready и /ping.
func NewRouter() *Router {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		RespondJSON(w, http.StatusOK, "pong")
	})

	return &Router{
		mux: mux,
	}
}

// InitRouter возвращает http.Handler, который обслуживает все зарегистрированные маршруты
// под префиксом /api/v{apiVer}.
func (m *Router) InitRouter(apiVer int) http.Handler {
	router := http.NewServeMux()

	pattern := fmt.Sprintf("/api/v%d/", apiVer)
	prefix := fmt.Sprintf("/api/v%d", apiVer)

	router.Handle(pattern, http.StripPrefix(prefix, m.mux))

	return router
}

// AddHandler регистрирует обработчик для заданного шаблона маршрута
// в базовом ServeMux.
func (m *Router) AddHandler(pattern string, handler http.HandlerFunc) {
	m.mux.HandleFunc(pattern, handler)
}
