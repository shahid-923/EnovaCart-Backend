package service

import (
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/dto"
	"ecommerce-app/internal/helper"
	"ecommerce-app/internal/repository"
	"errors"

	"ecommerce-app/config"
)

type CatalogueService struct {
	Repo   repository.CatalogueRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s CatalogueService) CreateCategory(input dto.CreateCategoryRequest) error {

	err := s.Repo.CreateCategory(&domain.Category{
	Name:         input.Name,
	ParentId:     input.ParentId,
	ImageUrl:     input.ImageUrl,
	DisplayOrder: input.DisplayOrder,

})
	return err
}

func (s CatalogueService) GetCategories() ([]*domain.Category, error) {

	categories, err := s.Repo.FindCategories()
	if err != nil {
		return nil, errors.New("categories dont exist")
	}

	return categories, nil
}

func (s CatalogueService) GetCategory(id int) (*domain.Category, error) {

	category, err := s.Repo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("category does not exist")
	}

	return category, nil
}

func (s CatalogueService) EditCategory(id int, input dto.CreateCategoryRequest) (*domain.Category, error) {

	existCat, err := s.Repo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("category does not exists")
	}

	if input.Name != "" {
		existCat.Name = input.Name
	}

	if input.ParentId > 0 {
		existCat.ParentId = input.ParentId
	}

	if input.ImageUrl != "" {
		existCat.ImageUrl = input.ImageUrl
	}

	if input.DisplayOrder > 0 {
		existCat.DisplayOrder = input.DisplayOrder
	}

	updatedCat, err := s.Repo.EditCategory(existCat)
	if err != nil {
		return nil, err
	}

	return updatedCat, nil
}

func (s CatalogueService) DeleteCategory(id int) error {
	err := s.Repo.DeleteCategory(id)

	if err != nil {
		return errors.New("category does not exist for deleting")
	}
	return nil
}

func (s CatalogueService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {
	err := s.Repo.CreateProduct(&domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryId:  input.CategoryId,
		ImageUrl:    input.ImageUrl,
		UserId:      int(user.ID),
		Stock:       uint(input.Stock),
	})

	if err != nil {
		return errors.New("failed to create product")
	}
	return err
}

func (s CatalogueService) GetProducts() ([]*domain.Product, error) {

	products, err := s.Repo.FindProducts()
	if err != nil {
		return nil, errors.New("products do not exist")
	}

	return products, nil
}

func (s CatalogueService) GetProductById(id int) (*domain.Product, error) {

	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return nil, errors.New("product does not exist")
	}

	return product, nil
}

func (s CatalogueService) EditProduct(id int, input dto.CreateProductRequest, user domain.User) (*domain.Product, error) {

	existingProduct, err := s.Repo.FindProductById(id)
	if err != nil {
		return nil, errors.New("product does not exist")
	}
    
	if existingProduct.UserId != int(user.ID){
		return nil, errors.New("you dont have manage rights of this product")
	}
	if input.Name != "" {
		existingProduct.Name = input.Name
	}

	if input.Description != "" {
		existingProduct.Description = input.Description
	}

	if input.Price > 0 {
		existingProduct.Price = input.Price
	}

	if input.CategoryId > 0 {
		existingProduct.CategoryId = input.CategoryId
	}

	if input.ImageUrl != "" {
		existingProduct.ImageUrl = input.ImageUrl
	}

	if input.Stock > 0 {
		existingProduct.Stock = input.Stock
	}

	updatedProduct, err := s.Repo.EditProduct(existingProduct)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}
func (s CatalogueService) DeleteProduct(id int, user domain.User) error {

	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return errors.New("product does not exist")
	}

	if product.UserId != int(user.ID) {
		return errors.New("you do not have permission to manage this product")
	}

	err = s.Repo.DeleteProduct(id)
	if err != nil {
		return err
	}

	return nil
}
func (s CatalogueService) GetSellerProducts(id int) ([]*domain.Product, error) {

	products, err := s.Repo.FindSellerProducts(id)
	if err != nil {
		return nil, errors.New("products do not exist")
	}

	return products, nil
}

func (s CatalogueService) UpdateProductStock(e domain.Product)(*domain.Product, error){

	product, err := s.Repo.FindProductById(int(e.ID))
	if err != nil{
		return nil, errors.New("product does not exist")
	}

	//verify product owner
	if product.UserId != e.UserId{    // if provided userid and db's product id are same then
		return nil, errors.New("you do not have permission to manage this product")
	}
	product.Stock = e.Stock

	editProduct, err := s.Repo.EditProduct(product)
	if err != nil{
		return nil, err
	}
	return editProduct, nil
}