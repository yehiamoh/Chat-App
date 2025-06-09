package database

import (
	"chat-app/pkg/models"
	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnection() *gorm.DB{
	err:=godotenv.Load()
	if err!=nil{
		fmt.Println("Error in Loading .env file")
	}
	/*dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	fmt.Print(dbHost,"\n",dbUser,"\n",dbPassword,"\n",dbName,"\n",dbPort)
	*/
	db, err := gorm.Open(postgres.Open("host=localhost user=chatuser password=pass dbname=chatapp port=5440 sslmode=disable"), &gorm.Config{})
	if err != nil {
		fmt.Println("error in connecting to the database : ",err)
		panic("Failed to connect to database")
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("Failed to migrate database")
	}
	return db

}