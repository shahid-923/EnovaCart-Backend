package handlers

import (
	"net/http"

	"ecommerce-app/internal/api/rest"
	"ecommerce-app/internal/dto"
	"ecommerce-app/internal/repository"
	"ecommerce-app/internal/service"

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

	pvtRoutes.Post("/cart", userHandler.AddToCart)
	pvtRoutes.Get("/cart", userHandler.GetCart)
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

func (h *UserHandler) GetProfile(ctx fiber.Ctx) error {
	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	profile, err := h.svc.GetProfile(user.ID)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(profile)
}

func (h *UserHandler) CreateProfile(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Profile created",
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
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Cart fetched",
	})
}

func (h *UserHandler) GetOrders(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Orders fetched",
	})
}

func (h *UserHandler) GetOrder(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":  "Order fetched",
		"order_id": id,
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
