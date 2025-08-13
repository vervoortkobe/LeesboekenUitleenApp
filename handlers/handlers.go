package handlers

import (
	"fmt"
	"log"
	"main/models"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Handler houdt een referentie naar de GORM-databaseverbinding vast.
type Handler struct {
	DB *gorm.DB
}

// --- Pagina Handlers ---

func (h *Handler) ShowKlaspagina(c *fiber.Ctx) error {
	zoekterm := c.Query("zoek")
	var klassen []models.Klas

	// Gebruik Preload om gerelateerde leerlingen direct mee te laden (eager loading).
	// De functie binnen Preload voegt een conditie toe aan de query voor leerlingen.
	query := h.DB.
		Preload("Leerlingen", "naam LIKE ?", "%"+zoekterm+"%").
		Where("naam LIKE ?", "%"+zoekterm+"%").
		Order("naam asc").
		Find(&klassen)

	if query.Error != nil {
		log.Printf("Error getting classes: %v", query.Error)
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.Render("klaspagina", fiber.Map{
		"Title":    "Klassenoverzicht",
		"Klassen":  klassen,
		"Zoekterm": zoekterm,
	}, "layouts/main")
}

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

// --- CRUD Handlers ---

func (h *Handler) AddKlas(c *fiber.Ctx) error {
	klas := models.Klas{Naam: c.FormValue("naam")}
	if klas.Naam != "" {
		h.DB.Create(&klas)
	}
	return c.Redirect("/")
}

func (h *Handler) DeleteKlas(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	// GORM's Select("Leerlingen") zorgt voor een cascade delete
	// van de gerelateerde leerlingen als de database dit ondersteunt.
	// Voor SQLite is het beter om dit handmatig te doen voor de zekerheid.
	h.DB.Where("klas_id = ?", id).Delete(&models.Leerling{})
	h.DB.Delete(&models.Klas{}, id)
	return c.Redirect("/")
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
	return c.Redirect("/")
}

func (h *Handler) DeleteLeerling(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	h.DB.Delete(&models.Leerling{}, id)
	return c.Redirect("/")
}

func (h *Handler) UpdateLeerling(c *fiber.Ctx) error {
	leerlingID, _ := c.ParamsInt("id")

	// Naam van leerling bijwerken
	h.DB.Model(&models.Leerling{}).Where("id = ?", leerlingID).Update("naam", c.FormValue("naam"))

	// Leesdata bijwerken
	form, _ := c.MultipartForm()
	for key, values := range form.Value {
		if strings.HasPrefix(key, "datum_") {
			boekIDStr := strings.TrimPrefix(key, "datum_")
			boekID, _ := c.ParamsInt(boekIDStr)
			datum := values[0]

			lb := models.LeerlingBoek{
				LeerlingID: uint(leerlingID),
				BoekID:     uint(boekID),
				Leesdatum:  datum,
			}
			// Clauses.OnConflict doet een "UPSERT": UPDATE als het record bestaat, INSERT als het nieuw is.
			h.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "leerling_id"}, {Name: "boek_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"leesdatum"}),
			}).Create(&lb)
		}
	}
	return c.Redirect(fmt.Sprintf("/leerling/%d", leerlingID))
}

func (h *Handler) AddBoek(c *fiber.Ctx) error {
	boek := models.Boek{
		Titel:     c.FormValue("titel"),
		AviNiveau: c.FormValue("avi_niveau"),
	}
	if boek.Titel != "" && boek.AviNiveau != "" {
		h.DB.Create(&boek)
	}
	return c.Redirect("/boeken")
}

func (h *Handler) DeleteBoek(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	h.DB.Delete(&models.Boek{}, id)
	return c.Redirect("/boeken")
}
