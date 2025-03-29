package handlers

import (
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type Repository struct {
	DB    *gorm.DB
	Model interface{}
}

func NewRepository(db *gorm.DB, model interface{}) *Repository {
	return &Repository{
		DB:    db,
		Model: model,
	}
}

// Create - Crée un nouvel enregistrement
func (r *Repository) Create(data interface{}) error {
	result := r.DB.Create(data)
	if result.Error != nil {
		return fmt.Errorf("erreur création: %w", result.Error)
	}
	return nil
}

// GetByIDOrUsername - Récupère par ID ou username
func (r *Repository) Get(identifier interface{}) (interface{}, error) {
	modelType := reflect.TypeOf(r.Model)
	newModel := reflect.New(modelType).Interface()

	query := r.DB.Model(r.Model)

	switch v := identifier.(type) {
	case int, uint:
		query = query.Where("id = ?", v)
	case string:
		query = query.Where("username = ?", v)
	default:
		return nil, errors.New("type d'identifiant non supporté")
	}

	if err := query.First(newModel).Error; err != nil {
		return nil, fmt.Errorf("non trouvé: %w", err)
	}

	return reflect.ValueOf(newModel).Elem().Interface(), nil
}

// GetAll - Récupère tous les enregistrements
func (r *Repository) GetAll() (interface{}, error) {
	modelType := reflect.TypeOf(r.Model)
	sliceType := reflect.SliceOf(modelType)
	results := reflect.New(sliceType).Interface()

	if err := r.DB.Find(results).Error; err != nil {
		return nil, fmt.Errorf("erreur récupération: %w", err)
	}

	return reflect.ValueOf(results).Elem().Interface(), nil
}

// Update - Met à jour un enregistrement
func (r *Repository) Update(identifier interface{}, updates map[string]interface{}) error {
	query := r.DB.Model(r.Model)

	switch v := identifier.(type) {
	case int, uint:
		query = query.Where("id = ?", v)
	case string:
		query = query.Where("username = ?", v)
	default:
		return errors.New("type d'identifiant non supporté")
	}

	if err := query.Updates(updates).Error; err != nil {
		return fmt.Errorf("erreur mise à jour: %w", err)
	}

	return nil
}

// Delete - Supprime un enregistrement
func (r *Repository) Delete(identifier interface{}) error {
	query := r.DB

	switch v := identifier.(type) {
	case int, uint:
		query = query.Where("id = ?", v)
	case string:
		query = query.Where("username = ?", v)
	default:
		return errors.New("type d'identifiant non supporté")
	}

	if err := query.Delete(r.Model).Error; err != nil {
		return fmt.Errorf("erreur suppression: %w", err)
	}

	return nil
}

func (r *Repository) GetByField(field string, value interface{}) (interface{}, error) {
	// Utilisation
	// user, err := userRepo.GetByField("email", "john@example.com")
	modelType := reflect.TypeOf(r.Model)
	newModel := reflect.New(modelType).Interface()

	if err := r.DB.Where(fmt.Sprintf("%s = ?", field), value).First(newModel).Error; err != nil {
		return nil, err
	}

	return reflect.ValueOf(newModel).Elem().Interface(), nil
}
