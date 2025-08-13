package main

import (
	"log"
	"main/handlers"
	"main/models"
	"os/exec"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

func main() {
	// Database connectie
	db, err := gorm.Open(sqlite.Open("db/leesapp.db"), &gorm.Config{})
	if err != nil {
		panic("DB connectie mislukt")
	}
	DB = db

	// Migreer modellen
	DB.AutoMigrate(&models.Klas{}, &models.Leerling{}, &models.AVINiveau{}, &models.Boek{}, &models.LeesDatum{})

	// Fiber setup met templates
	engine := html.New("./public", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Static files (voor CSS/JS)
	app.Static("/static", "./static")

	// Routes
	app.Get("/klassen", handlers.KlassenPagina)
	app.Post("/klassen", handlers.VoegKlasToe)
	app.Put("/klassen/:id", handlers.PasKlasAan)
	app.Delete("/klassen/:id", handlers.VerwijderKlas)

	app.Post("/leerlingen", handlers.VoegLeerlingToe)
	app.Put("/leerlingen/:id", handlers.PasLeerlingAan)
	app.Delete("/leerlingen/:id", handlers.VerwijderLeerling)

	app.Get("/boeken", handlers.BoekenPagina)
	app.Post("/niveaus", handlers.VoegNiveauToe)
	app.Post("/boeken", handlers.VoegBoekToe)
	app.Put("/boeken/:id", handlers.PasBoekAan)
	app.Delete("/boeken/:id", handlers.VerwijderBoek)

	app.Get("/leerlingen/:id", handlers.LeerlingDetail)
	app.Put("/leerlingen/:id/naam", handlers.UpdateLeerlingNaam)
	app.Post("/leerlingen/:id/leesdata", handlers.VoegLeesDatumToe)
	app.Put("/leesdata/:id", handlers.PasLeesDatumAan)

	// Start server en open browser
	go func() {
		// Wacht kort om te zorgen dat de server draait
		<-time.After(500 * time.Millisecond)
		openBrowser("http://localhost:3000/klassen")
	}()

	// Start server
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Server starten mislukt: %v", err)
	}
}
