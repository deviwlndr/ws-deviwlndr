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

	page.Get("/checkip", controller.Homepage) //ujicoba panggil package musik

	page.Get("/mahasiswa", controller.GetMahasiswa)

	page.Get("/mahasiswa/:id", controller.GetMahasiswaID) //menampilkan data presensi berdasarkan id

	// page.Post("/insert", controller.InsertDataPresensi)

	// page.Put("/update/:id", controller.UpdateData)

	// page.Delete("/delete/:id", controller.DeletePresensiByID)

	page.Get("/docs/*", swagger.HandlerDefault)
}
