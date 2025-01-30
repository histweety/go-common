package handlers

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SeedingData struct {
}

func Seeding(db *mongo.Database, path string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cwd, err := os.Getwd()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		file, err := os.Open(filepath.Join(cwd, path))
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
