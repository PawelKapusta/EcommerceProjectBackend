package database

import (
	"backend/database/models"
	"backend/utils"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math/rand"
)

var database *gorm.DB = nil

func Connect() {
	db, err := gorm.Open(sqlite.Open("database.db"))
	if err != nil {
		panic("Database not working")
	}

	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Company{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Payment{})
	db.AutoMigrate(&models.Item{})
	db.AutoMigrate(&models.Order{})
	db.AutoMigrate(&models.OrderProduct{})

	database = db
}

func GetDatabase() *gorm.DB {
	return database
}

func addDefaultProductsToDatabase(name string, price float32, categoryId uint, companyId uint, imageUrl string) {
	c := new(models.Product)
	err := database.Where("ID = ?", c.ID).Take(&c)
	if err.RowsAffected == 0 {
		database.Create(&models.Product{
			Name:        name,
			Price:       price,
			Description: utils.RandString(int(rand.Int63n(64))),
			CategoryID:  categoryId,
			CompanyID:   companyId,
			ImageUrl:    imageUrl,
		})
	}
}

func addDefaultCategoriesToDatabase(name string) {
	c := new(models.Category)
	err := database.Where("name = ?", name).Take(&c)
	if err.RowsAffected == 0 {
		database.Create(&models.Category{
			Name:        name,
			Description: utils.RandString(8),
		})
	}
}

func addDefaultCompaniesToDatabase(name string) {
	c := new(models.Company)
	err := database.Where("name = ?", name).Take(&c)
	if err.RowsAffected == 0 {
		database.Create(&models.Company{
			Name:        name,
			Description: utils.RandString(8),
		})
	}
}

func addAdminToDatabase() {
	c := new(models.User)
	err := database.Where("email = ?", "pawel@admin").Take(&c)
	if err.RowsAffected == 0 {
		database.Create(&models.User{Name: "Pawel", Surname: "Kapusta",
			Username: "sony", Email: "pawel@admin", Service: "Application", Token: "sony:" + uuid.NewString(), GoToken: uuid.NewString(), Password: "admin", Admin: true})
	}
}

func addDefaultUsersToDatabase(name string, email string) {
	c := new(models.User)
	err := database.Where("email = ?", email).Take(&c)
	if err.RowsAffected == 0 {
		database.Create(&models.User{Name: name, Surname: name,
			Username: name, Email: email, Service: "Application",
			Token: name + ":" + uuid.NewString(), GoToken: uuid.NewString(), Password: name + "test"})
	}

}

func InitDefaultDatabase() {
	productPrice := 1000.00
	allProducts := 1
	addAdminToDatabase()
	for i := 1; i <= 5; i++ {
		number := fmt.Sprint(i)
		addDefaultCategoriesToDatabase("Category" + number)
		addDefaultCompaniesToDatabase("Company" + number)
		for index := 1; index <= rand.Intn(10)+10; index++ {
			totalProductsSting := fmt.Sprint(allProducts)
			addDefaultProductsToDatabase("Product"+totalProductsSting, float32(productPrice), uint(i), uint(i), "https://softwaremathematics.com/wp-content/uploads/2021/02/ultimate-guide-to-your-product-launch.jpg")
			allProducts++
		}
		addDefaultUsersToDatabase("test"+number, "user"+number+"@app.com")
	}
}
