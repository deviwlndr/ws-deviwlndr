package url

import (
	"github.com/deviwlndr/ws-deviwlndr/controller"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/swagger" // swagger handler
)

func Web(page *fiber.App) {
	// page.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)  //API from user whatsapp message from iteung gowa
	// page.Get("/ws/whatsauth/qr", websocket.New(controller.WsWhatsAuthQR)) //websocket whatsauth

	page.Get("/", controller.Sink)
	page.Post("/", controller.Sink)
	page.Put("/", controller.Sink)
	page.Patch("/", controller.Sink)
	page.Delete("/", controller.Sink)
	page.Options("/", controller.Sink)

	//url mahasiswa
	page.Get("/checkip", controller.Homepage)
	page.Get("/mahasiswa", controller.GetMahasiswa)
	page.Get("/mahasiswa/:id", controller.GetMahasiswaID)
	page.Get("/mahasiswa/npm/:npm", controller.GetMahasiswaFromNPM)
	page.Get("/mahasiswa/npm/:id", controller.GetMahasiswaFromNPM)
	page.Post("/insertmahasiswa", controller.InsertDataMahasiswa)
	page.Put("/mahasiswa/update/:npm", controller.UpdateDataMahasiswa)
	page.Delete("/mahasiswa/delete/:npm", controller.DeleteMahasiswaByNPM)

	//url dosen
	page.Get("/dosen", controller.GetDosen)
	page.Get("/dosen/:kode_dosen", controller.GetDosenFromKodeDosen)
	page.Post("/insertdosen", controller.InsertDosen)

	page.Get("/docs/*", swagger.HandlerDefault)
}
