package main

import (
	"log"

	"github.com/deviwlndr/ws-deviwlndr/config"

	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2/middleware/cors"


	"github.com/deviwlndr/ws-deviwlndr/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(musik.Dangdut()))
}
