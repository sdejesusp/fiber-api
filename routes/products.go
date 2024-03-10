package routes

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sdejesusp/fiber-api/database"
	"github.com/sdejesusp/fiber-api/models"
)

type Product struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
	Price        string `json:"price"`
}

func CreateResponseProduct(productModel models.Product) Product {
	return Product{ID: productModel.ID, Name: productModel.Name, SerialNumber: productModel.SerialNumber, Price: productModel.Price}
}

func ValidateNewProduct(product models.Product) error {
	if err := ValidateProductPrice(product.Price); err != nil {
		return err
	}

	return nil
}

func ValidateProductPrice(productPrice string) error {
	price, err := strconv.ParseFloat(productPrice, 64)
	if err != nil {
		return errors.New("invaid product price")
	}

	if price < 0 {
		return errors.New("the price cannot be negative")
	}
	return nil
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := ValidateNewProduct(product); err != nil {
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

func FindProduct(id int, product *models.Product) error {
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

	if err := FindProduct(id, &product); err != nil {
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

	if err := FindProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateProduct struct {
		Name         string `json:"name"`
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

func ChangeProductPrice(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("please ensure that id is an integer")
	}

	var product models.Product
	if err := FindProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdatePrice struct {
		Price string `json:"price"`
	}

	var updateData UpdatePrice

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := ValidateProductPrice(updateData.Price); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	product.Price = updateData.Price
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

	if err := FindProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&product).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("product deleted")

}
