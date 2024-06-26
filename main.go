package main

import (
	"log"
	"github.com/deviwlndr/ws-deviwlndr/config"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/deviwlndr/ws-deviwlndr/url"
	
	"github.com/gofiber/fiber/v2"

	// docs are generated by Swag CLI, you have to import them.
	// replace with your own docs folder, usually "github.com/username/reponame/docs"
	_ "github.com/deviwlndr/ws-deviwlndr/docs"
)

// @title TES SWAGGER ULBI
// @version 1.0
// @description This is a sample swagger for Fiber

// @contact.name API Support
// @contact.url https://github.com/deviwlndr
// @contact.email 714220054@std.ulbi.ac.id

// @host ws-deviwlndr-59b3a0157dd7.herokuapp.com
// @BasePath /
// @schemes https http


func main() {
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(musik.Dangdut()))
}
