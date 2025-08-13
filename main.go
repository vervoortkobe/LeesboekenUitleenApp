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

	// Initialiseer de database met GORM
	db := database.InitDatabase("./boekendata.db")

	// Maak de handler aan met de GORM db-verbinding
	h := &handlers.Handler{DB: db}

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())

	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(staticFS),
		PathPrefix: "static",
	}))

	// --- Routes (deze blijven hetzelfde) ---
	app.Get("/", h.ShowKlaspagina)

	app.Post("/klas/toevoegen", h.AddKlas)
	app.Post("/klas/verwijderen/:id", h.DeleteKlas)

	app.Post("/leerling/toevoegen/:klas_id", h.AddLeerling)
	app.Post("/leerling/verwijderen/:id", h.DeleteLeerling)
	app.Get("/leerling/:id", h.ShowLeerlingpagina)
	app.Post("/leerling/aanpassen/:id", h.UpdateLeerling)

	app.Get("/boeken", h.ShowBoekenpagina)
	app.Post("/boek/toevoegen", h.AddBoek)
	app.Post("/boek/verwijderen/:id", h.DeleteBoek)

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
