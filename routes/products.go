package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sdejesusp/fiber-api/database"
	"github.com/sdejesusp/fiber-api/models"
)

type Product struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateResponseProduct(productModel models.Product) Product {
	return Product{ID: productModel.ID, Name: productModel.Name, SerialNumber: productModel.SerialNumber}
}


func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&product)

	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}



func GetProducts(c *fiber.Ctx) error {
	Products := []models.Product{}

	database.Database.Db.Find(&Products)
	responseProducts := []Product{}
	
	for _, product := range Products {
		responseProduct := CreateResponseProduct(product)
		responseProducts = append(responseProducts, responseProduct)
	}

	return c.Status(200).JSON(responseProducts)
}


func findProduct(id int, product *models.Product) error {
	database.Database.Db.Find(&product, "id=?", id)
	if product.ID == 0 {
		return errors.New("Product does not exist")
	}
	return nil
}


func GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("please ensure that id is an integer")
	}

	var product models.Product

	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseProduct := CreateResponseProduct(product)
	return c.Status(200).JSON(responseProduct)

}



func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var product models.Product

	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	
	type UpdateProduct struct {
		Name string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}

	var updateData UpdateProduct

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	product.Name = updateData.Name
	product.SerialNumber = updateData.SerialNumber

	database.Database.Db.Save(&product)
	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}


func DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var product models.Product

	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&product).Error; err != nil {
		return c.Status(404).JSON(err.Error() )
	}

	return c.Status(200).SendString("product deleted")

}