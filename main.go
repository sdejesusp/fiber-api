package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sdejesusp/fiber-api/database"
	"github.com/sdejesusp/fiber-api/routes"
)


func main() {
	database.ConnectDb()
	app := fiber.New()

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}


func welcome(c *fiber.Ctx) error {
	return c.SendString("hola")
}

func setupRoutes(app *fiber.App) {
	// welcome endpoint
	app.Get("/api", welcome)

	// User endponts
	users := app.Group("/api/users")
	users.Post("", routes.CreateUser)
	users.Get("", routes.GetUsers)
	users.Get("/:id", routes.GetUser)
	users.Put("/:id", routes.UpdateUser)
	users.Delete("/:id", routes.DeleteUser)

	// Product endpoints 
	products := app.Group("/api/products")
	products.Post("", routes.CreateProduct)
	products.Get("", routes.GetProducts)
	products.Get("/:id", routes.GetProduct)
	products.Put("/:id", routes.UpdateProduct)
	products.Delete("/:id", routes.DeleteProduct)

	// Order endpoints
	orders := app.Group("/api/orders")
	orders.Post("", routes.CreateOrder)
	orders.Get("", routes.GetOrders)
	orders.Get("/:id", routes.GetOrder)
}