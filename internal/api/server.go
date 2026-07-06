package api

import (
	"log"

	"ecommerce-app/config"
	"ecommerce-app/internal/api/rest"
	"ecommerce-app/internal/api/rest/handlers"
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/helper"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig) {
	log.Printf("Starting server")

	app := fiber.New()
	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})  //db is the gateway for doing database operations


	if err != nil {
		log.Fatalf("Database connection error: %v\n", err )
	}
    
	log.Println("Database connected successfully")
	
	//run migration
	err = db.AutoMigrate(&domain.User{}, &domain.BankAccount{}, &domain.Category{}, &domain.Product{})
	if err != nil {
		log.Fatalf("Database migration error: %v\n", err)
	}
    log.Println("migration was successfull")

    //cors configuration

    c := cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:3030"},
    AllowHeaders: []string{"Content-Type", "Accept", "Authorization"},    // prevents sql injections
    AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"},
    })
	
    app.Use(c)

	auth := helper.SetupAuth(config.AppSecret)

	rh := &rest.RestHandler{
		App: app,
		DB: db,
		Auth: auth,
		Config: config,
	}

	setupRoutes(rh)

	log.Println("Listening on port", config.ServerPort)

    if err := app.Listen(":" + config.ServerPort); err != nil {
	log.Fatalf("Server failed to start: %v", err)
   }
}


func setupRoutes(rh *rest.RestHandler) {
	handlers.SetupUserRoutes(rh)          // for userRoutes
	handlers.SetupCatalogueRoutes(rh)     // for catalogueRoutes
}