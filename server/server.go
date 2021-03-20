package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql/handler"
	"miikka.xyz/rs/config"
	"miikka.xyz/rs/graph/generated"
	"miikka.xyz/rs/graph/resolvers"
	"miikka.xyz/rs/oauth"
	"miikka.xyz/rs/server/handlers"
	"miikka.xyz/rs/store"
)

// Server ...
type Server struct {
	db     *store.DB
	config *config.Config
}

// NewServer gives you server instance
func NewServer(db *store.DB, config *config.Config) *Server {
	return &Server{
		db:     db,
		config: config,
	}
}

// Start server
func (s *Server) Start() {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers.NewRootResolvers(s.db)}))

	// Github authentication
	ghAuth := oauth.GetGithub(s.db)

	// Don't use prefix here. Routes under prefix are possibly under proxy
	http.HandleFunc("/login/github", ghAuth.GetStartHandler(s.config))
	http.HandleFunc("/github", ghAuth.GetCallbackHandler(s.config))

	// Arduino route: Exchange rfid to username
	http.HandleFunc(s.config.APIPrefix+s.config.DeviceTagRoute, handlers.TagToUsername(s.db, s.config))
	// Arduino route: Save game
	http.HandleFunc(s.config.APIPrefix+s.config.DeviceSaveRoute, handlers.SaveGame(s.db, s.config))
	// Basic health check for monitoring
	http.HandleFunc(s.config.APIPrefix+"/health", health)

	// Local development gives an error if CORS is not enabled
	if s.config.IsProduction {
		http.Handle(s.config.APIPrefix+"/graphql", srv)
	} else {
		http.Handle(s.config.APIPrefix+"/graphql", applyCORS(srv))
	}

	// Uncomment below if u want GraphQL plugin for browser
	// http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	env := "development"
	if s.config.IsProduction {
		env = "production"
	}
	log.Info("Server started in ", env, " mode")
	log.Info("Listening ", s.config.Port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+s.config.Port, nil))
}

// Set CORS-headers
func applyCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Health check
func health(w http.ResponseWriter, r *http.Request) {
	log.Info("health check")
	w.Write([]byte("ok"))
}
