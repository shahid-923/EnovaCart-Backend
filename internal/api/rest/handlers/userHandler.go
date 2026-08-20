package handlers

import (
	"ecommerce-app/internal/api/rest"
	"ecommerce-app/internal/dto"
	"ecommerce-app/internal/repository"
	"ecommerce-app/internal/service"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	svc *service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App

	svc := &service.UserService{
		Repo:   repository.NewUserRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
		CRepo:  repository.NewCatalogueRepository(rh.DB),
	}

	userHandler := &UserHandler{
		svc: svc,
	}

	pubRoutes := app.Group("/users")
	pubRoutes.Post("/register", userHandler.Register)
	pubRoutes.Post("/login", userHandler.Login)

	pvtRoutes := pubRoutes.Group("/", rh.Auth.Authorize())

	pvtRoutes.Get("/verify", userHandler.GetVerificationCode) // generates code
	pvtRoutes.Post("/verify", userHandler.Verify)             // verifies code

	pvtRoutes.Post("/profile", userHandler.CreateProfile)
	pvtRoutes.Get("/profile", userHandler.GetProfile)
	pvtRoutes.Patch("/profile", userHandler.UpdateProfile)

	pvtRoutes.Post("/cart", userHandler.AddToCart)
	pvtRoutes.Get("/cart", userHandler.GetCart)

	pvtRoutes.Post("/order", userHandler.CreateOrder)
	pvtRoutes.Get("/order", userHandler.GetOrders)
	pvtRoutes.Get("/order/:id", userHandler.GetOrder)

	pvtRoutes.Post("/become-seller", userHandler.BecomeSeller)
}

func (h *UserHandler) Register(ctx fiber.Ctx) error {
	var input dto.UserSignup

	if err := ctx.Bind().JSON(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
			"error":   err.Error(),
		})
	}

	_, token, err := h.svc.Register(input)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create user",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
		"token":   token,
	})
}

func (h *UserHandler) Login(ctx fiber.Ctx) error {
	var input dto.UserLogin

	if err := ctx.Bind().JSON(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
			"error":   err.Error(),
		})
	}

	token, err := h.svc.Login(input.Email, input.Password)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid credentials",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "login successful",
		"token":   token,
	})
}

func (h *UserHandler) CreateProfile(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	req := dto.ProfileInput{}

	if err = ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide a valid input",
		})
	}

	log.Printf("User %v, user", user)

	// create profile
	err = h.svc.CreateProfile(user.ID, req)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to create profile",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Profile created successfully",
	})
}

func (h *UserHandler) GetProfile(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	log.Println(user)

	//call user service and get profile
	profile, err := h.svc.GetProfile(user.ID)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to get profile",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Profile fetched",
		"profile": profile,
	})
}

func (h *UserHandler) UpdateProfile(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	req := dto.ProfileInput{}

	if err = ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide a valid input",
		})
	}

	err = h.svc.UpdateProfile(user.ID, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to update profile",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

func (h *UserHandler) Verify(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var req dto.VerificationCodeInput
	if err := ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide a valid input",
		})
	}

	err = h.svc.VerifyCode(user.ID, req.Code)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "code verified",
	})
}
func (h *UserHandler) GetVerificationCode(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	err = h.svc.GetVerificationCode(user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Verification email sent",
	})
}

func (h *UserHandler) AddToCart(ctx fiber.Ctx) error {

	req := dto.CreateCartRequest{}
	err := ctx.Bind().JSON(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide a valid product and quantity",
		})
	}

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	// cal userService and perform creat cart operation
	cartItems, err := h.svc.CreateCart(req, user)

	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "cart created succesfully", cartItems)
}

func (h *UserHandler) GetCart(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	cart, totalAmount, err := h.svc.FindCart(user.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":      "cart fetched successfully",
		"cart":         cart,
		"total_amount": totalAmount,
	})
}

func (h *UserHandler) CreateOrder(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return errors.New("user unauthenticated")
	}

	orederRef, err := h.svc.CreateOrder(user)

	if err != nil {
		return rest.InternalError(ctx, errors.New("unable to create order"))
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Order created successfully",
		"order":   orederRef,
	})
}

func (h *UserHandler) GetOrders(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return errors.New("user unauthenticated")
	}

	orders, err := h.svc.GetOrders(user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Orders fetched",
		"orders":  orders,
	})
}

func (h *UserHandler) GetOrder(ctx fiber.Ctx) error {

	orderId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid order id",
		})
	}

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	order, err := h.svc.GetOrderById(uint(orderId), user.ID)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "get order by id",
		"order":   order,
	})
}

func (h *UserHandler) BecomeSeller(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	req := dto.SellerInput{}

	err = ctx.Bind().JSON(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "request parameters are not valid",
		})
	}
	token, err := h.svc.BecomeSeller(user.ID, req)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to become seller",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Seller application submitted",
		"token":   token,
	})
}
