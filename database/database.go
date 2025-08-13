package database

import (
	"log"
	"main/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase initialiseert de GORM-databaseverbinding en migreert de schema's.
func InitDatabase(filepath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log alle GORM queries
	})
	if err != nil {
		log.Fatalf("Fout bij het verbinden met de database: %v", err)
	}

	log.Println("Databaseverbinding succesvol.")

	// AutoMigrate zal de tabellen aanmaken/bijwerken op basis van de modellen.
	// Het zal geen kolommen of tabellen verwijderen.
	err = db.AutoMigrate(&models.Klas{}, &models.Leerling{}, &models.Boek{}, &models.LeerlingBoek{})
	if err != nil {
		log.Fatalf("Fout bij het migreren van de database: %v", err)
	}

	log.Println("Database succesvol gemigreerd.")
	return db
}
