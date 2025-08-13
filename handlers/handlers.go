package handlers

import (
	"fmt"
	"log"
	"main/models"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	DB *gorm.DB
}

// ShowKlaspagina blijft ongewijzigd
func (h *Handler) ShowKlaspagina(c *fiber.Ctx) error {
	zoekterm := c.Query("zoek")
	var klassen []models.Klas

	preloadLeerlingen := func(db *gorm.DB) *gorm.DB {
		return db.Order("naam asc").Where("naam LIKE ?", "%"+zoekterm+"%")
	}

	dbQuery := h.DB.Preload("Leerlingen", preloadLeerlingen).Order("naam asc")

	var gefilterdeKlassen []models.Klas
	if err := dbQuery.Find(&klassen).Error; err != nil {
		log.Printf("Error getting classes: %v", err)
		return c.Status(500).SendString("Internal Server Error")
	}

	if zoekterm != "" {
		for _, klas := range klassen {
			if strings.Contains(strings.ToLower(klas.Naam), strings.ToLower(zoekterm)) || len(klas.Leerlingen) > 0 {
				gefilterdeKlassen = append(gefilterdeKlassen, klas)
			}
		}
	} else {
		gefilterdeKlassen = klassen
	}

	return c.Render("klaspagina", fiber.Map{
		"Title":    "Klassenoverzicht",
		"Klassen":  gefilterdeKlassen,
		"Zoekterm": zoekterm,
	}, "layouts/main")
}

// ShowBoekenpagina blijft ongewijzigd
func (h *Handler) ShowBoekenpagina(c *fiber.Ctx) error {
	zoekterm := c.Query("zoek")
	var boeken []models.Boek

	h.DB.Where("titel LIKE ? OR avi_niveau LIKE ?", "%"+zoekterm+"%", "%"+zoekterm+"%").
		Order("avi_niveau, titel").
		Find(&boeken)

	grouped := make(map[string][]models.Boek)
	var aviNiveaus []string
	for _, b := range boeken {
		if _, ok := grouped[b.AviNiveau]; !ok {
			aviNiveaus = append(aviNiveaus, b.AviNiveau)
		}
		grouped[b.AviNiveau] = append(grouped[b.AviNiveau], b)
	}
	sort.Strings(aviNiveaus)

	var result []models.AviGroep
	for _, niveau := range aviNiveaus {
		result = append(result, models.AviGroep{Niveau: niveau, Boeken: grouped[niveau]})
	}

	return c.Render("boekenpagina", fiber.Map{
		"Title":      "Boekenoverzicht",
		"AviGroepen": result,
		"Zoekterm":   zoekterm,
	}, "layouts/main")
}

// ShowLeerlingpagina blijft ongewijzigd
func (h *Handler) ShowLeerlingpagina(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var leerling models.Leerling
	if err := h.DB.First(&leerling, id).Error; err != nil {
		return c.Status(404).SendString("Leerling niet gevonden")
	}

	var alleBoeken []models.Boek
	h.DB.Order("avi_niveau, titel").Find(&alleBoeken)

	var gelezenData []models.LeerlingBoek
	h.DB.Where("leerling_id = ?", id).Find(&gelezenData)

	gelezenMap := make(map[uint]string)
	for _, data := range gelezenData {
		gelezenMap[data.BoekID] = data.Leesdatum
	}

	grouped := make(map[string][]models.BoekMetDatum)
	var aviNiveaus []string
	for _, boek := range alleBoeken {
		bmd := models.BoekMetDatum{Boek: boek, Leesdatum: gelezenMap[boek.ID]}
		if _, ok := grouped[boek.AviNiveau]; !ok {
			aviNiveaus = append(aviNiveaus, boek.AviNiveau)
		}
		grouped[boek.AviNiveau] = append(grouped[boek.AviNiveau], bmd)
	}
	sort.Strings(aviNiveaus)

	return c.Render("leerlingpagina", fiber.Map{
		"Title":         fmt.Sprintf("Details voor %s", leerling.Naam),
		"Leerling":      leerling,
		"GroupedBoeken": grouped,
		"AviNiveaus":    aviNiveaus,
	}, "layouts/main")
}

// --- CRUD Handlers met RedirectToRoute ---

func (h *Handler) AddKlas(c *fiber.Ctx) error {
	klas := models.Klas{Naam: c.FormValue("naam")}
	if klas.Naam != "" {
		h.DB.Create(&klas)
	}
	return c.RedirectToRoute("klassen.index", fiber.Map{}) // CORRECTIE
}

func (h *Handler) DeleteKlas(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	h.DB.Where("klas_id = ?", id).Delete(&models.Leerling{})
	h.DB.Delete(&models.Klas{}, id)
	return c.RedirectToRoute("klassen.index", fiber.Map{}) // CORRECTIE
}

func (h *Handler) AddLeerling(c *fiber.Ctx) error {
	klasID, _ := c.ParamsInt("klas_id")
	leerling := models.Leerling{
		Naam:   c.FormValue("naam"),
		KlasID: uint(klasID),
	}
	if leerling.Naam != "" {
		h.DB.Create(&leerling)
	}
	return c.RedirectToRoute("klassen.index", fiber.Map{}) // CORRECTIE
}

func (h *Handler) DeleteLeerling(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	h.DB.Delete(&models.Leerling{}, id)
	return c.RedirectToRoute("klassen.index", fiber.Map{}) // CORRECTIE
}

func (h *Handler) UpdateLeerling(c *fiber.Ctx) error {
	leerlingID, _ := c.ParamsInt("id")

	naam := c.FormValue("naam")
	if naam != "" {
		h.DB.Model(&models.Leerling{}).Where("id = ?", leerlingID).Update("naam", naam)
	}

	c.Request().PostArgs().VisitAll(func(keyBytes, valueBytes []byte) {
		key := string(keyBytes)
		if strings.HasPrefix(key, "datum_") {
			boekIDStr := strings.TrimPrefix(key, "datum_")
			boekID, err := strconv.Atoi(boekIDStr)
			if err != nil {
				return
			}
			datum := string(valueBytes)

			if datum != "" {
				lb := models.LeerlingBoek{LeerlingID: uint(leerlingID), BoekID: uint(boekID), Leesdatum: datum}
				h.DB.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "leerling_id"}, {Name: "boek_id"}},
					DoUpdates: clause.AssignmentColumns([]string{"leesdatum"}),
				}).Create(&lb)
			} else {
				h.DB.Where("leerling_id = ? AND boek_id = ?", leerlingID, boekID).Delete(&models.LeerlingBoek{})
			}
		}
	})

	// Stuur terug naar de detailpagina van de leerling, met de ID als parameter
	return c.RedirectToRoute("leerling.show", fiber.Map{"id": leerlingID}) // CORRECTIE
}

func (h *Handler) AddBoek(c *fiber.Ctx) error {
	boek := models.Boek{
		Titel:     c.FormValue("titel"),
		AviNiveau: c.FormValue("avi_niveau"),
	}
	if boek.Titel != "" && boek.AviNiveau != "" {
		h.DB.Create(&boek)
	}
	return c.RedirectToRoute("boeken.index", fiber.Map{}) // CORRECTIE
}

func (h *Handler) DeleteBoek(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	h.DB.Delete(&models.Boek{}, id)
	return c.RedirectToRoute("boeken.index", fiber.Map{}) // CORRECTIE
}

func (h *Handler) UpdateBoek(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	titel := c.FormValue("titel")
	aviNiveau := c.FormValue("avi_niveau")

	// Validatie: zorg ervoor dat de velden niet leeg zijn
	if titel != "" && aviNiveau != "" {
		h.DB.Model(&models.Boek{}).Where("id = ?", id).Updates(models.Boek{
			Titel:     titel,
			AviNiveau: aviNiveau,
		})
	}

	return c.RedirectToRoute("boeken.index", fiber.Map{})
}
