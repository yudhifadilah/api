package practice

import (
	"context"
	"net/http"
	"practice-api/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func Index(c *fiber.Ctx) error {
	var notes []models.Latihan
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil catatan",
		})
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var note models.Latihan
		if err := cursor.Decode(&note); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal mengambil catatan",
			})
		}
		notes = append(notes, note)
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"notes": notes})
}

func Show(c *fiber.Ctx) error {
	id := c.Params("id")
	var note models.Latihan
	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Data Tidak Ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil catatan",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"notes": note})
}

func Create(c *fiber.Ctx) error {
	var note models.Latihan
	if err := c.BodyParser(&note); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Request Body Invalid",
		})
	}
	_, err := collection.InsertOne(context.TODO(), note)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat catatan",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"notes": note})
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var note models.Latihan
	if err := c.BodyParser(&note); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Request Body Invalid",
		})
	}
	_, err := collection.ReplaceOne(context.TODO(), bson.M{"id": id}, note)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate catatan",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Catatan berhasil diupdate"})
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus catatan",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Catatan berhasil dihapus"})
}
