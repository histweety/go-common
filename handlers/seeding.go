package handlers

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Seeding(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cwd, err := os.Getwd()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		dirPath := filepath.Join(cwd, "seeders")
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if filepath.Ext(file.Name()) != ".json" {
				continue
			}

			if fileHasExecuted(db, file.Name()) {
				continue
			}

			dataFile, err := os.Open(filepath.Join(dirPath, file.Name()))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			defer dataFile.Close()

			if err := seedData(db, dataFile); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	}
}

func seedData(db *mongo.Database, file *os.File) error {
	var seeders map[string]interface{}
	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &seeders)

	for collectionName, collectionValues := range seeders {
		collection := db.Collection(collectionName)
		_, err := collection.InsertMany(context.Background(), collectionValues.([]interface{}))
		if err != nil {
			return err
		}
	}

	migration := db.Collection("migrations")
	migration.InsertOne(context.Background(), map[string]interface{}{
		"filename": file.Name(),
	})

	return nil
}

func fileHasExecuted(db *mongo.Database, filename string) bool {
	migration := db.Collection("migrations")
	count, _ := migration.CountDocuments(context.Background(), map[string]interface{}{
		"filename": filename,
	})

	return count > 0
}
