package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kkovac1/products/service/products"
	"github.com/kkovac1/products/service/user"
)

type ApiServer struct {
	address string
	db      *sql.DB
}

func NewApiServer(address string, db *sql.DB) *ApiServer {
	return &ApiServer{
		address: address,
		db:      db,
	}
}

func (server *ApiServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(server.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productsStore := products.NewStore(server.db)
	productsHandler := products.NewHandler(productsStore)
	productsHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", server.address)
	return http.ListenAndServe(server.address, router)
}
