package controller

import (
	"errors"
	"fmt"
	"strconv" 
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

func GetMahasiswaFromNPM(c *fiber.Ctx) error {
    npm := c.Params("npm") // mengambil parameter npm dari URL
    if npm == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "NPM parameter is required",
        })
    }

    // Mengonversi npm yang berupa string ke integer
    npmInt, err := strconv.Atoi(npm)  // Mengonversi npm string menjadi integer
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Invalid NPM format, should be a number",
        })
    }

    // Panggil fungsi GetMahasiswaFromNPM untuk mengambil data berdasarkan npm
    mahasiswa := cek.GetMahasiswaFromNPM(npmInt)  // Panggil fungsi dengan npm dalam format integer

    // Periksa apakah data mahasiswa ditemukan
    if mahasiswa.Npm == 0 {  // Jika Npm 0 berarti data tidak ditemukan
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
	var mahasiswa inimodel.Mahasiswa  // Pastikan menggunakan model yang tepat
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
		mahasiswa.Npm,         // NPM mahasiswa untuk mencari data
		mahasiswa,             // Struct Mahasiswa yang berisi data yang ingin diupdate
	)

	// Jika terjadi error atau update gagal
	if err != nil {
		// Menyediakan pesan error yang lebih jelas jika update gagal
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Failed to update mahasiswa data: %v", err),
		})
	}

	if !success {
		// Jika tidak ada data yang diupdate (misalnya mahasiswa tidak ditemukan)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": fmt.Sprintf("No mahasiswa found with NPM %d", mahasiswa.Npm),
		})
	}

	// Jika berhasil, kirim respon sukses
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Mahasiswa data successfully updated",
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
    err = cek.DeleteMahasiswaByNPM(npmInt)  // Pastikan fungsi ini ada di dalam cek
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

//DOSEN 
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
		"status": http.StatusOK,
		"message": fmt.Sprintf("Data found for kode_dosen %d", kodeDosenInt),
		"dosen": dosen,
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
		"status":   http.StatusOK,
		"message":  "Dosen successfully inserted",
		"insertedID": insertedID,
	})
}

// UpdateDataDosen handles the HTTP request for updating a dosen record based on kode_dosen
func UpdateDataDosen(c *fiber.Ctx) error {
	// Get the kode_dosen from the URL parameter
	kodeDosen := c.Params("kode_dosen")

	// Parse kode_dosen into an integer
	kodeDosenInt, err := strconv.Atoi(kodeDosen)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid kode_dosen format",
		})
	}

	// Parse the request body into a Dosen object
	var dosen inimodel.Dosen
	if err := c.BodyParser(&dosen); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Failed to parse request body",
		})
	}

	// Log the kode_dosen and check if it exists
	fmt.Printf("Attempting to update dosen with kode_dosen: %d\n", kodeDosenInt)

	// Call the UpdateDosen function with kode_dosen and the Dosen object
	updated, err := cek.UpdateDosen(kodeDosenInt, dosen)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Check if update was successful
	if !updated {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": fmt.Sprintf("No dosen found with kode_dosen %d", kodeDosenInt),
		})
	}

	// Return success response
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Data successfully updated",
	})
}






func DeleteDosenByKodeDosen(c *fiber.Ctx) error {
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

	// Panggil fungsi DeleteDosenByKodeDosen untuk menghapus data dosen
	err = cek.DeleteDosenByKodeDosen(kodeDosenInt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Error deleting data for kode_dosen %d", kodeDosenInt),
			"error":   err.Error(),
		})
	}

	// Jika berhasil, kirim respon sukses
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Dosen with kode_dosen %d successfully deleted", kodeDosenInt),
	})
}

