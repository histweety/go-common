package handlers

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	var seeders map[string][]map[string]interface{}
	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &seeders)

	for collectionName, collectionValues := range seeders {
		collection := db.Collection(collectionName)

		for _, values := range collectionValues {
			// Support for relational data
			relations := catchHasRelations(values)
			if len(relations) > 0 {
				for _, relation := range relations {
					parts := strings.Split(relation, ":")
					valueKey, searchKey, collection := parts[0], parts[1], parts[2]
					relatedCollection := db.Collection(collection)
					var actualValue map[string]interface{}
					relatedCollection.FindOne(context.Background(), bson.M{searchKey: values[relation]}).Decode(&actualValue)
					values[valueKey] = actualValue["_id"]
					delete(values, relation)
				}
			}
		}

		_, err := collection.InsertMany(context.Background(), collectionValues)
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

func catchHasRelations(values map[string]interface{}) []string {
	relations := []string{}

	for key := range values {
		parts := strings.Split(key, ":")
		if len(parts) == 3 {
			relations = append(relations, key)
		}
	}

	return relations
}
