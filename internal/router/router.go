package router

import (
	"fmt"
	"net/http"

	handler "github.com/Quiqui-dev/Earthbound/internal/api/handler/user"
	"github.com/Quiqui-dev/Earthbound/middleware"
	"github.com/Quiqui-dev/Earthbound/util"
)

func NewRouter(userH *handler.UserHandler) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", userH.CreateUser)
	mux.HandleFunc("/users/login", userH.Login)
	mux.HandleFunc("/users/logout", userH.Logout)

	chain := &middleware.Chain{}
	chain.Use(middleware.RecoverMiddleware)
	chain.Use(middleware.CorsMiddleware)

	wrappedMux := chain.Then(mux)

	mainMux := http.NewServeMux()

	mainMux.Handle("/api/", http.StripPrefix("/api", wrappedMux))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", util.GetEnv("PORT", "8080")),
		Handler: mainMux,
	}

	return srv
}
