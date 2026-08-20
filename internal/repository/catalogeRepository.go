package repository

import (
	"errors"
	"log"

	"ecommerce-app/internal/domain"

	"gorm.io/gorm"
)

type CatalogueRepository interface {
	// Categories
	CreateCategory(c *domain.Category) error
	FindCategories() ([]*domain.Category, error)
	FindCategoryById(id int) (*domain.Category, error)
	EditCategory(c *domain.Category) (*domain.Category, error)
	DeleteCategory(id int) error

	// Products
	CreateProduct(c *domain.Product) error
	FindProducts() ([]*domain.Product, error)
	FindProductById(id int) (*domain.Product, error)
	FindSellerProducts(id int) ([]*domain.Product, error)
	EditProduct(c *domain.Product) (*domain.Product, error)
	DeleteProduct(id int) error
}
type catalogueRepository struct { // actual implementation and contains db connection
	db *gorm.DB
}

func NewCatalogueRepository(db *gorm.DB) CatalogueRepository { // constructor creates new instance returns as interface type to keep code loosely coupled
	return &catalogueRepository{
		db: db,
	}
}

func (c catalogueRepository) CreateCategory(e *domain.Category) error {

	err := c.db.Create(e).Error // Creates a new record by automatically inferring the model from 'e'.

	if err != nil {
		log.Printf("db_err: %v", err)
		return errors.New("creating category failed")
	}

	return nil
}

func (c catalogueRepository) FindCategories() ([]*domain.Category, error) {

	var categories []*domain.Category

	err := c.db.Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c catalogueRepository) FindCategoryById(id int) (*domain.Category, error) {

	var category domain.Category

	err := c.db.First(&category, id).Error
	if err != nil {
		log.Printf("db_err: %v", err)
		return nil, errors.New("category does not exist")
	}

	return &category, nil
}

func (c catalogueRepository) EditCategory(e *domain.Category) (*domain.Category, error) {

	err := c.db.Save(e).Error

	if err != nil {
		log.Printf("db_err: %v", err)
		return nil, errors.New("failed to update category")
	}

	return e, nil
}

func (c catalogueRepository) DeleteCategory(id int) error {

	err := c.db.Delete(&domain.Category{}, id).Error

	if err != nil {
		log.Printf("db_err: %v", err)
		return errors.New("failed to delete category")
	}

	return nil
}

func (c catalogueRepository) CreateProduct(e *domain.Product) error {

	err := c.db.Model(&domain.Product{}).Create(e).Error // Explicitly specifies the Product model before creating the record.
	if err != nil {
		log.Printf("db_err: %v", err)
		return errors.New("cannot create product")
	}

	return nil
}

func (c catalogueRepository) FindProducts() ([]*domain.Product, error) {

	var products []*domain.Product

	err := c.db.Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
func (c catalogueRepository) FindProductById(id int) (*domain.Product, error) {
	var product domain.Product

	err := c.db.First(&product, id).Error
	if err != nil {
		log.Printf("db_err: %v", err)
		return nil, errors.New("product does not exist")
	}

	return &product, nil
}

func (c catalogueRepository) EditProduct(e *domain.Product) (*domain.Product, error) {

	err := c.db.Save(e).Error

	if err != nil {
		log.Printf("db_err: %v", err)
		return nil, errors.New("failed to update product")
	}

	return e, nil
}

func (c catalogueRepository) DeleteProduct(id int) error {

	result := c.db.Delete(&domain.Product{}, id)

	if result.Error != nil {
		log.Printf("db_err: %v", result.Error)
		return errors.New("failed to delete product")
	}

	if result.RowsAffected == 0 {
		return errors.New("product does not exist")
	}

	return nil
}

func (c catalogueRepository) FindSellerProducts(id int) ([]*domain.Product, error) {

	var products []*domain.Product

	err := c.db.Where("user_id = ?", id).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
