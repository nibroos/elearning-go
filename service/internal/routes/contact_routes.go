package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

func SetupContactRoutes(contacts fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	contactRepo := repository.NewContactRepository(gormDB, sqlDB)
	contactService := service.NewContactService(contactRepo)
	contactController := rest.NewContactController(contactService)

	// prefix /contacts

	contacts.Post("/index-contact", contactController.ListContacts)
	contacts.Post("/show-contact", contactController.GetContactByID)
	contacts.Post("/create-contact", contactController.CreateContact)
	contacts.Post("/update-contact", contactController.UpdateContact)
	contacts.Post("/delete-contact", contactController.DeleteContact)
	contacts.Post("/restore-contact", contactController.RestoreContact)
	contacts.Post("/auth-index-contact", contactController.ListContactsByAuthUser)
	contacts.Post("/auth-show-contact", contactController.GetContactByIDByAuthUser)
	contacts.Post("/auth-create-contact", contactController.CreateContactByAuthUser)
	contacts.Post("/auth-update-contact", contactController.UpdateContactByAuthUser)
	contacts.Post("/auth-delete-contact", contactController.DeleteContactByAuthUser)
	contacts.Post("/auth-restore-contact", contactController.RestoreContactByAuthUser)
}
