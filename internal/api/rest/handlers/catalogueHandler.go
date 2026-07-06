package handlers

import (
	"ecommerce-app/internal/api/rest"
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/dto"
	"ecommerce-app/internal/repository"
	"ecommerce-app/internal/service"

	"strconv"

	"github.com/gofiber/fiber/v3"
)

type CatalogueHandler struct {
	svc *service.CatalogueService
}

func SetupCatalogueRoutes(rh *rest.RestHandler) {
	app := rh.App

	svc := &service.CatalogueService{
		Repo:   repository.NewCatalogueRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}

	handler := &CatalogueHandler{
		svc: svc,
	}

	// Public Routes
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProduct)

	app.Get("/categories", handler.GetCategories)

	// Seller Routes
	sellRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller())

	// Categories
	sellRoutes.Get("/categories/:id", handler.GetCategoryById)
	sellRoutes.Post("/categories", handler.CreateCategories)
	sellRoutes.Patch("/categories/:id", handler.EditCategory)
	sellRoutes.Delete("/categories/:id", handler.DeleteCategory)

	// Products
	sellRoutes.Post("/products", handler.CreateProducts)
	sellRoutes.Get("/products", handler.GetProducts)
	sellRoutes.Get("/products/:id", handler.GetProduct)
	sellRoutes.Put("/products/:id", handler.EditProduct)
	sellRoutes.Patch("/products/:id", handler.UpdateStock)
	sellRoutes.Delete("/products/:id", handler.DeleteProduct)
	sellRoutes.Get("/my-products", handler.FindSellerProducts)
}

// ===== Category Handlers =====
func (h *CatalogueHandler) CreateCategories(ctx fiber.Ctx) error {

	req := dto.CreateCategoryRequest{}

	err := ctx.Bind().JSON(&req)

	if err != nil {
		return rest.BadRequestError(ctx, "create category is not valid")
	}

	err = h.svc.CreateCategory(req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "category created successfully", nil)

}

func (h *CatalogueHandler) GetCategories(ctx fiber.Ctx) error {

	cats, err := h.svc.GetCategories()
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "categories fetched successfully", cats)
}
func (h *CatalogueHandler) GetCategoryById(ctx fiber.Ctx) error {

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "invalid category id")
	}

	cat, err := h.svc.GetCategory(id)
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessResponse(ctx, "category fetched successfully", cat)
}

func (h *CatalogueHandler) EditCategory(ctx fiber.Ctx) error {

	req := dto.CreateCategoryRequest{}
	id, _ := strconv.Atoi(ctx.Params("id"))

	err := ctx.Bind().JSON(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "update category request is not valid")
	}

	updated, err := h.svc.EditCategory(id, req)

	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "category edited successfully", updated)

}

func (h *CatalogueHandler) DeleteCategory(ctx fiber.Ctx) error {

	id, _ := strconv.Atoi(ctx.Params("id"))
	err := h.svc.DeleteCategory(id)

	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "category deleted successfully", nil)
}

// ===== Product Handlers =====

func (h *CatalogueHandler) CreateProducts(ctx fiber.Ctx) error {

	req := dto.CreateProductRequest{}
	err := ctx.Bind().JSON(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "create product request is not valid")
	}

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	err = h.svc.CreateProduct(req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "product created successfully", nil)
}

func (h *CatalogueHandler) GetProducts(ctx fiber.Ctx) error {
	
	products, err := h.svc.GetProducts()
	if err != nil{
		return rest.ErrorMessage(ctx, 404, err)
	}
    
	return rest.SuccessResponse(ctx, "products fetched successfully", products)
}

func (h *CatalogueHandler) GetProduct(ctx fiber.Ctx) error {
	
	id, _ := strconv.Atoi(ctx.Params("id"))
	product, err := h.svc.GetProductById(id)

	if err != nil{
	return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessResponse(ctx, "product fetched successfully", product)
}

func (h *CatalogueHandler) EditProduct(ctx fiber.Ctx) error {

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "invalid product id")
	}

	req := dto.CreateProductRequest{}

	err = ctx.Bind().JSON(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "edit product request is not valid")
	}

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	product, err := h.svc.EditProduct(id, req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "product updated successfully", product)
}

func (h *CatalogueHandler) UpdateStock(ctx fiber.Ctx) error {

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "invalid product id")
	}

	req := dto.UpdateStockRequest{}

	err = ctx.Bind().JSON(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "update stock request is not valid")
	}

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	product := domain.Product{
		ID: uint(id),
		Stock: uint(req.Stock),
		UserId: int(user.ID),
	}
	
	updatedProduct, err := h.svc.UpdateProductStock(product)
	if err != nil {
	return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "stock updated successfully", updatedProduct)
}

func (h *CatalogueHandler) DeleteProduct(ctx fiber.Ctx) error {

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "invalid product id")
	}

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	err = h.svc.DeleteProduct(id, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "product deleted successfully", nil)
}

func (h *CatalogueHandler) FindSellerProducts(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	products, err := h.svc.GetSellerProducts(int(user.ID))
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "seller products fetched successfully", products)
}