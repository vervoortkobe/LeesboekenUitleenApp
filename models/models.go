package models

import (
	"time"

	"gorm.io/gorm"
)

type Klas struct {
	gorm.Model
	Naam       string
	Leerlingen []Leerling `gorm:"foreignKey:KlasID"`
}

type Leerling struct {
	gorm.Model
	Naam     string
	KlasID   uint
	LeesData []LeesDatum `gorm:"foreignKey:LeerlingID"`
}

type AVINiveau struct {
	gorm.Model
	Niveau string
	Boeken []Boek `gorm:"foreignKey:AVINiveauID"`
}

type Boek struct {
	gorm.Model
	Titel       string
	AVINiveauID uint
}

type LeesDatum struct {
	gorm.Model
	LeerlingID uint
	BoekID     uint
	Datum      *time.Time
}
