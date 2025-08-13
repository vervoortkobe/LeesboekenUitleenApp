package models

import "gorm.io/gorm"

// Klas vertegenwoordigt een klas. Een klas heeft meerdere leerlingen.
type Klas struct {
	gorm.Model
	Naam       string     `gorm:"unique;not null"`
	Leerlingen []Leerling `gorm:"foreignKey:KlasID"` // Has-Many relatie
}

// Leerling vertegenwoordigt een leerling. Een leerling behoort tot één klas.
type Leerling struct {
	gorm.Model
	Naam   string `gorm:"not null"`
	KlasID uint   // Foreign key voor de Klas
}

// Boek vertegenwoordigt een boek in de bibliotheek.
type Boek struct {
	gorm.Model
	Titel     string `gorm:"not null"`
	AviNiveau string `gorm:"not null"`
}

// LeerlingBoek is de 'join table' voor de many-to-many relatie
// tussen Leerlingen en Boeken. We definiëren deze expliciet
// omdat we een extra veld `Leesdatum` willen opslaan.
type LeerlingBoek struct {
	LeerlingID uint `gorm:"primaryKey"`
	BoekID     uint `gorm:"primaryKey"`
	Leesdatum  string
}

// --- Helper structs voor de templates ---

// BoekMetDatum is voor de leerlingpagina.
type BoekMetDatum struct {
	Boek
	Leesdatum string
}

// AviGroep is voor de boekenpagina.
type AviGroep struct {
	Niveau string
	Boeken []Boek
}
