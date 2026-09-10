package web

import (
	"html/template"
	"net/http"
)

type Server struct {
	mux			*http.ServeMux
	templates	*template.Template
}

func NewServer() (*Server, error) {
	templates, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}

	s := &Server{
		mux:		http.NewServeMux(),
		templates:	templates,
	}

	s.routes()

	return s, nil
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.Handle(
		"GET /static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("web/dist")),
		),
	)

	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /_ui/ping", s.handlePing)
	s.mux.HandleFunc("GET /", s.handleHome)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := s.templates.ExecuteTemplate(w, "home.html", nil); err != nil {
		http.Error(
			w, "failed to render page",
			http.StatusInternalServerError,
		)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`)) // Utter nonsense.
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(
		`<span class="status">htmx is alive</span>`,
	))
}
