package handlers

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type seed_dir struct {
	Filename string `json:"file_name"`
}

func Seeding(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := new(seed_dir)
		if err := c.BodyParser(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": 400,
				"message":    "Error parsing body",
				"data":       err.Error(),
			})
		}

		cwd, err := os.Getwd()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		file, err := os.Open(filepath.Join(cwd, payload.Filename))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		defer file.Close()

		var seeders map[string]interface{}
		byteValue, _ := io.ReadAll(file)
		json.Unmarshal(byteValue, &seeders)

		for collectionName, collectionValues := range seeders {
			collection := db.Collection(collectionName)
			_, err := collection.InsertMany(c.Context(), collectionValues.([]interface{}))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	}
}
