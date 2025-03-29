package main

import (
	"TransactoR/database"
	"TransactoR/middleware"

	// "TransactoR/router"
	"TransactoR/router"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"uniqueIndex"`
}

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
	if err := database.AutoMigrate(db, []interface{}{&User{}}); err != nil {
		log.Fatal(err)
	}

	// Création routeur
	r := router.New(db)

	// Déclaration route
	r.AddRoute("/users", "POST", createUser)

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

	user := User{Name: "John", Email: "john@example.com"}
	if err := tx.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
