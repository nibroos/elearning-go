package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

func SetupAddressRoutes(addresses fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	addressRepo := repository.NewAddressRepository(gormDB, sqlDB)
	addressService := service.NewAddressService(addressRepo)
	addressController := rest.NewAddressController(addressService)

	// prefix /addresses

	addresses.Post("/index-address", addressController.ListAddresses)
	addresses.Post("/show-address", addressController.GetAddressByID)
	addresses.Post("/create-address", addressController.CreateAddress)
	addresses.Post("/update-address", addressController.UpdateAddress)
	addresses.Post("/delete-address", addressController.DeleteAddress)
	addresses.Post("/restore-address", addressController.RestoreAddress)
	addresses.Post("/auth-index-address", addressController.ListAddressesByAuthUser)
	addresses.Post("/auth-show-address", addressController.GetAddressByIDByAuthUser)
	addresses.Post("/auth-create-address", addressController.CreateAddressByAuthUser)
	addresses.Post("/auth-update-address", addressController.UpdateAddressByAuthUser)
	addresses.Post("/auth-delete-address", addressController.DeleteAddressByAuthUser)
	addresses.Post("/auth-restore-address", addressController.RestoreAddressByAuthUser)
}
