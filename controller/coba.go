package controller

import (
	"errors"
	"fmt"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2"
	inimodel"github.com/mhrndiva/kemahasiswaan/model"
	cek "github.com/mhrndiva/kemahasiswaan/module"
	"github.com/deviwlndr/ws-deviwlndr/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

)

// GetMahasiswaID godoc
// @Summary Get By ID Data Presensi.
// @Description Ambil per ID data presensi.
// @Tags Presensi
// @Accept json
// @Produce json
// @Param id path string true "Masukan ID"
// @Success 200 {object} Presensi
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /presensi/{id} [get]
func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func GetMahasiswa(c *fiber.Ctx) error {
	ps := cek.GetAllMahasiswa()
	return c.JSON(ps)
	
}

func GetMahasiswaID(c *fiber.Ctx) error {
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
	ps, err := cek.GetMahasiswaFromID(objID, config.Ulbimongoconn, "presensi")
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

func InsertDataMahasiswa(c *fiber.Ctx) error {
	var mahasiswa inimodel.Mahasiswa

	// Parsing body request ke struct Mahasiswa
	if err := c.BodyParser(&mahasiswa); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Pemanggilan fungsi cek.InsertMahasiswa dengan parameter sesuai
	insertedID := cek.InsertMahasiswa(
		mahasiswa.Nama,           // Nama mahasiswa
		mahasiswa.Phone_number,   // Phone number mahasiswa
		mahasiswa.Jurusan,        // Jurusan mahasiswa
		mahasiswa.Npm,            // NPM mahasiswa
		mahasiswa.Alamat,         // Alamat mahasiswa
		mahasiswa.Email,          // Email mahasiswa
		mahasiswa.Poin,           // Poin mahasiswa
	)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Data berhasil disimpan.",
		"inserted_id": insertedID,
	})
}

func UpdateDataMahasiswa(c *fiber.Ctx) error {
	// Ambil NPM dari parameter URL
	npm := c.Params("npm")
	if npm == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "NPM is required",
		})
	}

	// Parsing body request ke struct Mahasiswa
	var mahasiswa inimodel.Mahasiswa
	if err := c.BodyParser(&mahasiswa); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Validasi data Mahasiswa
	if mahasiswa.Npm == 0 || mahasiswa.Nama == "" || mahasiswa.Phone_number == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "NPM, Nama, and Phone_number are required",
		})
	}

	// Panggil fungsi UpdateMahasiswa dengan NPM dan struct Mahasiswa
	success, err := cek.UpdateMahasiswa(
		npm,         // NPM mahasiswa untuk mencari data
		mahasiswa,   // Struct Mahasiswa yang berisi data yang ingin diupdate
	)

	// Jika terjadi error atau update gagal
	if err != nil || !success {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update mahasiswa data",
			"error":   err.Error(),
		})
	}

	// Jika berhasil, kirim respon sukses
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Mahasiswa data successfully updated",
	})
}



// InsertDataPresensi godoc
// @Summary Insert data presensi.
// @Description Input data presensi.
// @Tags Presensi
// @Accept json
// @Produce json
// @Param request body ReqPresensi true "Payload Body [RAW]"
// @Success 200 {object} Presensi
// @Failure 400
// @Failure 500
// @Router /insert [post]
// func InsertMahasiswa(c *fiber.Ctx) error {
// 	db := config.Ulbimongoconn
// 	var presensi inimodel.Mahasiswa


// 	if err := c.BodyParser(&presensi); err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": err.Error(),
// 		})
// 	}


// 	insertedID, err := cek.InsertMahasiswa(db, "mahasiswa",
// 		presensi.Nama,         // string
// 		presensi.Npm,          // string
// 		presensi.Phone_number, // string
// 		presensi.Poin,         // int
// 		presensi.Jurusan,      // string
// 		presensi.Email,        // string
// 	)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  http.StatusInternalServerError,
// 			"message": err.Error(),
// 		})
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"status":      http.StatusOK,
// 		"message":     "Data berhasil disimpan.",
// 		"inserted_id": insertedID,
// 	})
// }

// UpdateData godoc
// @Summary Update data presensi.
// @Description Ubah data presensi.
// @Tags Presensi
// @Accept json
// @Produce json
// @Param id path string true "Masukan ID"
// @Param request body ReqPresensi true "Payload Body [RAW]"
// @Success 200 {object} Presensi
// @Failure 400
// @Failure 500
// @Router /update/{id} [put]
	// func UpdateData(c *fiber.Ctx) error {
	// 	db := config.Ulbimongoconn
	
	// 	// Get the ID from the URL parameter
	// 	id := c.Params("id")
	
	// 	// Parse the ID into an ObjectID
	// 	objectID, err := primitive.ObjectIDFromHex(id)
	// 	if err != nil {
	// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 			"status":  http.StatusInternalServerError,
	// 			"message": err.Error(),
	// 		})
	// 	}
	
	// 	// Parse the request body into a Presensi object
	// 	var presensi inimodel.Presensi
	// 	if err := c.BodyParser(&presensi); err != nil {
	// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 			"status":  http.StatusInternalServerError,
	// 			"message": err.Error(),
	// 		})
	// 	}
	
	// 	// Call the UpdatePresensi function with the parsed ID and the Presensi object
	// 	err = cek.UpdatePresensi(db, "presensi",
	// 		objectID,
	// 		presensi.Longitude,
	// 		presensi.Latitude,
	// 		presensi.Location,
	// 		presensi.Phone_number,
	// 		presensi.Checkin,
	// 		presensi.Biodata)
	// 	if err != nil {
	// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 			"status":  http.StatusInternalServerError,
	// 			"message": err.Error(),
	// 		})
	// 	}
	
	// 	return c.Status(http.StatusOK).JSON(fiber.Map{
	// 		"status":  http.StatusOK,
	// 		"message": "Data successfully updated",
	// 	})
	// }

// DeletePresensiByID godoc
// @Summary Delete data presensi.
// @Description Hapus data presensi.
// @Tags Presensi
// @Accept json
// @Produce json
// @Param id path string true "Masukan ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /delete/{id} [delete]
	// func DeletePresensiByID(c *fiber.Ctx) error {
	// 	id := c.Params("id")
	// 	if id == "" {
	// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 			"status":  http.StatusInternalServerError,
	// 			"message": "Wrong parameter",
	// 		})
	// 	}
	
	// 	objID, err := primitive.ObjectIDFromHex(id)
	// 	if err != nil {
	// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
	// 			"status":  http.StatusBadRequest,
	// 			"message": "Invalid id parameter",
	// 		})
	// 	}
	
	// 	err = cek.DeletePresensiByID(objID, config.Ulbimongoconn, "presensi")
	// 	if err != nil {
	// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 			"status":  http.StatusInternalServerError,
	// 			"message": fmt.Sprintf("Error deleting data for id %s", id),
	// 		})
	// 	}
	
	// 	return c.Status(http.StatusOK).JSON(fiber.Map{
	// 		"status":  http.StatusOK,
	// 		"message": fmt.Sprintf("Data with id %s deleted successfully", id),
	// 	})
	// }	