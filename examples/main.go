package main

import (
	"TransactoR/database"
	"TransactoR/handlers"
	"TransactoR/middleware"
	"TransactoR/model"

	// "TransactoR/router"
	"TransactoR/router"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	// Configuration DB
	cfg := database.Config{
		Driver:          "postgres",
		DSN:             "host=localhost user=lincoln password=admin dbname=postgres port=5432 sslmode=disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
	}

	// Initialisation DB
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Migrations
	if err := database.AutoMigrate(db, []interface{}{&model.User{}}); err != nil {
		log.Fatal(err)
	}

	// Création routeur
	r := router.New(db)

	// Déclaration des routes
	r.AddRoute("/users", "POST", createUser)
	// r.AddRoute("/users", "GET", handlers.Get())
	// r.AddRoute("/users/all", "GET", handlers.GetUsersHandler)
	// r.AddRoute("/users", "PUT", handlers.UpdateUserHandler)
	// r.AddRoute("/users", "DELETE", handlers.DeleteUserHandler)

	// Création du repository générique
	userRepo := handlers.NewRepository(db, model.User{})

	// Création des handlers
	userHandler := handlers.NewCRUDHandler(userRepo)

	// Enregistrement des routes CRUD
	handlers.RCRUDRoutes(r.Router, "/users", userHandler)

	// Démarrage du serveur
	http.ListenAndServe(":8080", r)

	// Démarrage serveur
	log.Println("Serveur démarré sur :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	tx, err := middleware.GetTx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := model.User{Name: "John", Email: "john@example.com"}
	if err := tx.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
