package handlers

import (
	"main/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

func KlassenPagina(c *fiber.Ctx) error {
	var klassen []models.Klas
	DB.Preload("Leerlingen").Find(&klassen)
	return c.Render("klassen", fiber.Map{
		"Klassen": klassen,
	})
}

func VoegKlasToe(c *fiber.Ctx) error {
	klas := models.Klas{}
	if err := c.BodyParser(&klas); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	DB.Create(&klas)
	return c.JSON(klas)
}

func PasKlasAan(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	klas := models.Klas{Model: gorm.Model{ID: uint(id)}}
	if err := c.BodyParser(&klas); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	DB.Save(&klas)
	return c.JSON(klas)
}

func VerwijderKlas(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	DB.Delete(&models.Klas{}, id)
	return c.SendStatus(200)
}

func VoegLeerlingToe(c *fiber.Ctx) error {
	leerling := models.Leerling{}
	if err := c.BodyParser(&leerling); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	DB.Create(&leerling)
	return c.JSON(leerling)
}

func PasLeerlingAan(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	leerling := models.Leerling{Model: gorm.Model{ID: uint(id)}}
	if err := c.BodyParser(&leerling); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	DB.Save(&leerling)
	return c.JSON(leerling)
}

func VerwijderLeerling(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	DB.Delete(&models.Leerling{}, id)
	return c.SendStatus(200)
}

func BoekenPagina(c *fiber.Ctx) error {
	var niveaus []models.AVINiveau
	DB.Preload("Boeken").Find(&niveaus)
	return c.Render("boeken", fiber.Map{
		"Niveaus": niveaus,
	})
}

func VoegNiveauToe(c *fiber.Ctx) error {
	niveau := models.AVINiveau{}
	if err := c.BodyParser(&niveau); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	DB.Create(&niveau)
	return c.JSON(niveau)
}

func VoegBoekToe(c *fiber.Ctx) error {
	boek := models.Boek{}
	if err := c.BodyParser(&boek); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	DB.Create(&boek)
	return c.JSON(boek)
}

func PasBoekAan(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	boek := models.Boek{Model: gorm.Model{ID: uint(id)}}
	if err := c.BodyParser(&boek); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	DB.Save(&boek)
	return c.JSON(boek)
}

func VerwijderBoek(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	DB.Delete(&models.Boek{}, id)
	return c.SendStatus(200)
}

func LeerlingDetail(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var leerling models.Leerling
	DB.Preload(clause.Associations).First(&leerling, id)

	var niveaus []models.AVINiveau
	DB.Preload("Boeken").Find(&niveaus)

	leesMap := make(map[uint]*time.Time)
	for _, ld := range leerling.LeesData {
		leesMap[ld.BoekID] = ld.Datum
	}

	return c.Render("leerling", fiber.Map{
		"Leerling": leerling,
		"Niveaus":  niveaus,
		"LeesMap":  leesMap,
	})
}

func UpdateLeerlingNaam(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	leerling := models.Leerling{Model: gorm.Model{ID: uint(id)}}
	if err := c.BodyParser(&leerling); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	DB.Model(&leerling).Update("Naam", leerling.Naam)
	return c.JSON(leerling)
}

func VoegLeesDatumToe(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	type Req struct {
		BoekID uint   `json:"boekID"`
		Datum  string `json:"datum"`
	}
	req := Req{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	datum, err := time.Parse("2006-01-02", req.Datum)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Ongeldige datum"})
	}
	ld := models.LeesDatum{LeerlingID: uint(id), BoekID: req.BoekID, Datum: &datum}
	DB.Create(&ld)
	return c.JSON(ld)
}

func PasLeesDatumAan(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	type Req struct {
		Datum string `json:"datum"`
	}
	req := Req{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	datum, err := time.Parse("2006-01-02", req.Datum)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Ongeldige datum"})
	}
	DB.Model(&models.LeesDatum{Model: gorm.Model{ID: uint(id)}}).Update("Datum", &datum)
	return c.JSON(fiber.Map{"status": "success"})
}
