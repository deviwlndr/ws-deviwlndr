package controller

import (
	"errors"
	"fmt"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	cek "github.com/indrariksa/cobapakcage/module"
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