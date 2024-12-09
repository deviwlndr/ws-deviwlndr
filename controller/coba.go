package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aiteung/musik"
	"github.com/deviwlndr/ws-deviwlndr/config"
	"github.com/gofiber/fiber/v2"

	inimodel "github.com/mhrndiva/kemahasiswaan/model"
	cek "github.com/mhrndiva/kemahasiswaan/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
func GetDosen(c *fiber.Ctx) error {
	ps := cek.GetAllDosen()

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

func GetMahasiswaFromNPM(c *fiber.Ctx) error {
	npm := c.Params("npm") // mengambil parameter npm dari URL
	if npm == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "NPM parameter is required",
		})
	}

	// Mengonversi npm yang berupa string ke integer
	npmInt, err := strconv.Atoi(npm) // Mengonversi npm string menjadi integer
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid NPM format, should be a number",
		})
	}

	// Panggil fungsi GetMahasiswaFromNPM untuk mengambil data berdasarkan npm
	mahasiswa := cek.GetMahasiswaFromNPM(npmInt) // Panggil fungsi dengan npm dalam format integer

	// Periksa apakah data mahasiswa ditemukan
	if mahasiswa.Npm == 0 { // Jika Npm 0 berarti data tidak ditemukan
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": fmt.Sprintf("No data found for NPM %d", npmInt),
		})
	}

	// Mengembalikan data mahasiswa yang ditemukan
	return c.JSON(fiber.Map{
		"status":    http.StatusOK,
		"message":   "Data found",
		"mahasiswa": mahasiswa,
	})
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
		mahasiswa.Nama,         // Nama mahasiswa
		mahasiswa.Phone_number, // Phone number mahasiswa
		mahasiswa.Jurusan,      // Jurusan mahasiswa
		mahasiswa.Npm,          // NPM mahasiswa
		mahasiswa.Alamat,       // Alamat mahasiswa
		mahasiswa.Email,        // Email mahasiswa
		mahasiswa.Poin,         // Poin mahasiswa
	)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Data berhasil disimpan.",
		"inserted_id": insertedID,
	})
}

func UpdateDataMahasiswa(c *fiber.Ctx) error {
	// Ambil NPM dari URL parameter
	npmParam := c.Params("npm")

	/// Konversi NPM dari string ke integer
	npm, err := strconv.Atoi(npmParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid NPM format",
		})
	}

	// Parse body request ke struct Mahasiswa
	var mahasiswa inimodel.Mahasiswa
	if err := c.BodyParser(&mahasiswa); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
		})
	}

	// Panggil fungsi UpdateMahasiswa
	success, err := cek.UpdateMahasiswa(npm, mahasiswa)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Jika tidak ada mahasiswa yang diperbarui
	if !success {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": fmt.Sprintf("No mahasiswa found with NPM %d", npm),
		})
	}

	// Jika berhasil
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Data successfully updated",
	})
}

func DeleteMahasiswaByNPM(c *fiber.Ctx) error {
	npm := c.Params("npm")
	if npm == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "NPM is required",
		})
	}

	// Mengonversi npm menjadi integer
	npmInt, err := strconv.Atoi(npm)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid NPM format, should be a number",
		})
	}

	// Panggil fungsi untuk menghapus mahasiswa berdasarkan NPM
	err = cek.DeleteMahasiswaByNPM(npmInt) // Pastikan fungsi ini ada di dalam cek
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error deleting data for NPM %d", npmInt),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Data with NPM %d deleted successfully", npmInt),
	})
}

// DOSEN
func GetDosenFromKodeDosen(c *fiber.Ctx) error {
	// Ambil kode_dosen dari parameter URL
	kodeDosen := c.Params("kode_dosen")
	if kodeDosen == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "kode_dosen parameter is required",
		})
	}

	// Mengonversi kode_dosen menjadi integer
	kodeDosenInt, err := strconv.Atoi(kodeDosen)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid kode_dosen format, should be a number",
		})
	}

	// Panggil fungsi GetDosenFromKodeDosen untuk mendapatkan data dosen
	dosen := cek.GetDosenFromKodeDosen(kodeDosenInt)
	if dosen.Kode_dosen == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": fmt.Sprintf("No data found for kode_dosen %d", kodeDosenInt),
		})
	}

	// Jika berhasil, kirim respon dengan data dosen
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Data found for kode_dosen %d", kodeDosenInt),
		"dosen":   dosen,
	})
}

// InsertDosen handles the insertion of a new Dosen
func InsertDosen(c *fiber.Ctx) error {
	// Parsing body request ke struct Dosen
	var dosen inimodel.Dosen
	if err := c.BodyParser(&dosen); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Validasi data Dosen
	if dosen.Nama == "" || dosen.Kode_dosen == 0 || dosen.Phone_number == "" || dosen.Matkul == "" || dosen.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Nama, Kode_dosen, Phone_number, Matkul, and Email are required",
		})
	}

	// Panggil fungsi InsertDosen untuk menyimpan data dosen ke database
	insertedID := cek.InsertDosen(dosen.Nama, dosen.Kode_dosen, dosen.Phone_number, dosen.Matkul, dosen.Email)

	// Jika berhasil, kirim respon sukses dengan ID yang dimasukkan
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":     http.StatusOK,
		"message":    "Dosen successfully inserted",
		"insertedID": insertedID,
	})
}
