package webserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// WebServer holds configuration and routes for the HTTP server.
type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

// NewWebServer creates a new WebServer with sensible defaults.
func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

// AddHandler registers an HTTP handler for a given method+path pattern.
// Pattern format: "METHOD /path", e.g. "GET /order"
func (ws *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	ws.Handlers[method+" "+path] = handler
}

// Start sets up middleware, registers routes, and starts listening.
func (ws *WebServer) Start() {
	ws.Router.Use(middleware.Logger)
	ws.Router.Use(middleware.Recoverer)

	for pattern, handler := range ws.Handlers {
		var method, path string
		fmt.Sscanf(pattern, "%s %s", &method, &path)
		switch method {
		case "GET":
			ws.Router.Get(path, handler)
		case "POST":
			ws.Router.Post(path, handler)
		case "PUT":
			ws.Router.Put(path, handler)
		case "DELETE":
			ws.Router.Delete(path, handler)
		}
	}

	fmt.Printf("Starting web server on port %s\n", ws.WebServerPort)
	http.ListenAndServe(ws.WebServerPort, ws.Router)
}
