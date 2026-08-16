package rest

import (
	"ecommerce-app/config"
	"ecommerce-app/internal/helper"
	"ecommerce-app/pkg/payment"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type RestHandler struct { //holds all the dependencies your HTTP request handlers
	App    *fiber.App
	DB     *gorm.DB
	Auth   helper.Auth
	Config config.AppConfig
	Pc     payment.PaymentClient
}
