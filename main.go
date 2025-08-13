package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"main/database"
	"main/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

// openBrowser opent de standaardbrowser naar de opgegeven URL
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // Linux, BSD, etc.
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		log.Printf("Fout bij het openen van de browser: %v", err)
	}
}

//go:embed templates/*
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	templateRoot, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		log.Fatal(err)
	}

	engine := html.NewFileSystem(http.FS(templateRoot), ".html")
	engine.AddFunc("safe", func(s string) template.HTML {
		return template.HTML(s)
	})

	db := database.InitDatabase("./boekendata.db")
	h := &handlers.Handler{DB: db}
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(logger.New())

	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(staticFS),
		PathPrefix: "static",
	}))

	// --- Routes met Namen ---
	app.Get("/", func(c *fiber.Ctx) error {
		// Stuur direct door naar de genoemde route
		return c.RedirectToRoute("klassen.index", fiber.Map{})
	})

	// Klassenpagina en CRUD
	app.Get("/klassen", h.ShowKlaspagina).Name("klassen.index") // Geef de route een naam
	app.Post("/klas/toevoegen", h.AddKlas)
	app.Post("/klas/verwijderen/:id", h.DeleteKlas)

	// Leerling CRUD
	app.Post("/leerling/toevoegen/:klas_id", h.AddLeerling)
	app.Post("/leerling/verwijderen/:id", h.DeleteLeerling)
	app.Get("/leerling/:id", h.ShowLeerlingpagina).Name("leerling.show") // Geef de route een naam
	app.Post("/leerling/aanpassen/:id", h.UpdateLeerling)

	// Boekenpagina en CRUD
	app.Get("/boeken", h.ShowBoekenpagina).Name("boeken.index") // Geef de route een naam
	app.Post("/boek/toevoegen", h.AddBoek)
	app.Post("/boek/verwijderen/:id", h.DeleteBoek)
	app.Post("/boek/aanpassen/:id", h.UpdateBoek)

	// Start server en open browser
	go func() {
		// Wacht kort om te zorgen dat de server draait
		<-time.After(500 * time.Millisecond)
		openBrowser("http://localhost/")
	}()

	// Start server
	log.Println("Applicatie gestart op http://localhost:80")
	if err := app.Listen(":80"); err != nil {
		log.Fatalf("Server starten mislukt: %v", err)
	}
}
