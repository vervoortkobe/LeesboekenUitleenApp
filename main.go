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
	db, err := gorm.Open(sqlite.Open("db/boekendata.db"), &gorm.Config{})
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

	// Maak Handler instantie
	handler := handlers.NewHandler(DB)

	// Routes
	app.Get("/klassen", handler.KlassenPagina)
	app.Post("/klassen", handler.VoegKlasToe)
	app.Put("/klassen/:id", handler.PasKlasAan)
	app.Delete("/klassen/:id", handler.VerwijderKlas)

	app.Post("/leerlingen", handler.VoegLeerlingToe)
	app.Put("/leerlingen/:id", handler.PasLeerlingAan)
	app.Delete("/leerlingen/:id", handler.VerwijderLeerling)

	app.Get("/boeken", handler.BoekenPagina)
	app.Post("/niveaus", handler.VoegNiveauToe)
	app.Post("/boeken", handler.VoegBoekToe)
	app.Put("/boeken/:id", handler.PasBoekAan)
	app.Delete("/boeken/:id", handler.VerwijderBoek)

	app.Get("/leerlingen/:id", handler.LeerlingDetail)
	app.Put("/leerlingen/:id/naam", handler.UpdateLeerlingNaam)
	app.Post("/leerlingen/:id/leesdata", handler.VoegLeesDatumToe)
	app.Put("/leesdata/:id", handler.PasLeesDatumAan)

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
