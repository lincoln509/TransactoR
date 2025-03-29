package main

import (
	"TransactoR/database"
	"TransactoR/dbctx"
	"TransactoR/logging"
	"TransactoR/routes"
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm/logger"
)

func main() {
	// Initialisation DB
	cfg := database.Config{
		Driver:   "postgres",
		DSN:      "host=localhost user=postgres dbname=app port=5432",
		LogLevel: logger.Info,
	}
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Modèles à migrer
	// type User struct {
	// }

	// type Product struct {
	// }
	// models := []interface{}{
	// 	&User{},
	// 	&Product{},
	// }

	// Initialisation avec migrations
	// db, err := database.InitDBWithMigrations(cfg, models)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Création routeur
	logger := &logging.DefaultLogger{}
	router := routes.NewRouter(logger)

	// Déclaration route transactionnelle
	router.AddTransactionalRoute(routes.RouteConfig{
		Path:    "/users",
		Method:  "POST",
		Handler: createUserHandler,
	})

	db.http.ListenAndServe(":8080", router)
	// http.ListenAndServe(":8080", router)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	tx, err := dbctx.TxFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Logique métier
	user := User{Name: "Test"}
	if err := tx.Create(&user).Error; err != nil {
		http.Error(w, "Erreur de création", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
