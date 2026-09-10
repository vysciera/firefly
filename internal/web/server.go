package web

import (
	"html/template"
	"net/http"

	"firefly/internal/flowerpress"
)

type Server struct {
	mux         *http.ServeMux
	templates   *template.Template
	flowerpress *flowerpress.Client
}

type homePageData struct {
	User *flowerpress.User
}

func NewServer(flowerpressClient *flowerpress.Client) (*Server, error) {
	templates, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}

	s := &Server{
		mux:         http.NewServeMux(),
		templates:   templates,
		flowerpress: flowerpressClient,
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
			noCache(
				http.FileServer(http.Dir("web/dist")),
			),
		),
	)

	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /_ui/ping", s.handlePing)
	s.mux.HandleFunc("GET /_ui/flowerpress-health", s.handleFlowerpressHealth)
	s.mux.HandleFunc("GET /", s.handleHome)

	s.mux.HandleFunc("GET /login", s.handleLoginPage)
	s.mux.HandleFunc("POST /login", s.handleLogin)
	s.mux.HandleFunc("POST /logout", s.handleLogout)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := homePageData{}
	session := sessionFromRequest(r)

	if session != "" {
		user, cookies, err := s.flowerpress.Me(r.Context(), session)

		relayCookies(w, cookies)
		if err == nil {
			data.User = user
		}
	}

	if err := s.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(
			w, "failed to render page",
			http.StatusInternalServerError,
		)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`)) // Still nonsense.
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(
		`<span class="status">htmx is alive</span>`,
	))
}

func (s *Server) handleFlowerpressHealth(w http.ResponseWriter, r *http.Request) {
	health, err := s.flowerpress.Health(r.Context())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)

		_, _ = w.Write([]byte(
			`<span class="status">flowerpress unavailable</span>`,
		))

		return
	}

	if health.Status == "ok" {
		_, _ = w.Write([]byte(
			`<span class="status">flowerpress online</span>`,
		))

		return
	}

	_, _ = w.Write([]byte(
		`<span class="status">flowerpress unhealthy</span>`,
	))
}

func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
