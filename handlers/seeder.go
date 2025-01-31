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

func Seeder(DB *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var succeed []string = []string{}
		var skipped []string = []string{}

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
			if file.IsDir() ||
				filepath.Ext(file.Name()) != ".json" ||
				fileHasExecuted(DB, file.Name()) {
				skipped = append(skipped, file.Name())
				continue
			}

			dataFile, err := os.Open(filepath.Join(dirPath, file.Name()))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			defer dataFile.Close()

			if err := seedData(DB, dataFile); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			migration := DB.Collection("migrations")
			migration.InsertOne(context.Background(), map[string]interface{}{
				"filename": file.Name(),
			})
			succeed = append(succeed, file.Name())
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"statusCode": 200,
			"message":    "Seeding data success",
			"data":       map[string]interface{}{"succeed": succeed, "skipped": skipped},
		})
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

	return nil
}

func fileHasExecuted(db *mongo.Database, filename string) bool {
	migration := db.Collection("migrations")
	count, _ := migration.CountDocuments(context.Background(), map[string]interface{}{
		"filename": filename,
	})

	return count > 0
}
