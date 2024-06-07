package controller

import (
	"errors"
	"fmt"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	inimodel "github.com/deviwlndr/cobapackage/model"
	cek "github.com/deviwlndr/cobapackage/module"
	"github.com/deviwlndr/ws-deviwlndr/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func GetPresensi(c *fiber.Ctx) error {
	ps := cek.GetAllPresensi(config.Ulbimongoconn, "presensi")
	return c.JSON(ps)
	}

	func GetPresensiID(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "Wrong parameter",
			})
		}
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"status":  http.StatusBadRequest,
				"message": "Invalid id parameter",
			})
		}
		ps, err := cek.GetPresensiFromID(objID, config.Ulbimongoconn, "presensi")
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{
					"status":  http.StatusNotFound,
					"message": fmt.Sprintf("No data found for id %s", id),
				})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": fmt.Sprintf("Error retrieving data for id %s", id),
			})
		}
		return c.JSON(ps)
	}

	func InsertDataPresensi(c *fiber.Ctx) error {
		db := config.Ulbimongoconn
		var presensi inimodel.Presensi
		if err := c.BodyParser(&presensi); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
		insertedID, err := cek.InsertPresensi(db, "presensi",
			presensi.Longitude,
			presensi.Latitude,
			presensi.Location,
			presensi.Phone_number,
			presensi.Checkin,
			presensi.Biodata)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":      http.StatusOK,
			"message":     "Data berhasil disimpan.",
			"inserted_id": insertedID,
		})
	}

	func UpdateData(c *fiber.Ctx) error {
		db := config.Ulbimongoconn
	
		// Get the ID from the URL parameter
		id := c.Params("id")
	
		// Parse the ID into an ObjectID
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
	
		// Parse the request body into a Presensi object
		var presensi inimodel.Presensi
		if err := c.BodyParser(&presensi); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
	
		// Call the UpdatePresensi function with the parsed ID and the Presensi object
		err = cek.UpdatePresensi(db, "presensi",
			objectID,
			presensi.Longitude,
			presensi.Latitude,
			presensi.Location,
			presensi.Phone_number,
			presensi.Checkin,
			presensi.Biodata)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
	
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":  http.StatusOK,
			"message": "Data successfully updated",
		})
	}
	
	func DeletePresensiByID(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "Wrong parameter",
			})
		}
	
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"status":  http.StatusBadRequest,
				"message": "Invalid id parameter",
			})
		}
	
		err = cek.DeletePresensiByID(objID, config.Ulbimongoconn, "presensi")
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": fmt.Sprintf("Error deleting data for id %s", id),
			})
		}
	
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":  http.StatusOK,
			"message": fmt.Sprintf("Data with id %s deleted successfully", id),
		})
	}	