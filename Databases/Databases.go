package databases

import (
	"Backend-Projek-Aslab/Models/Oauth"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var UDB *gorm.DB

func ConnectUserAuth() {
	postgreSQL := os.Getenv("POSTGRESQL_URL")
	connectUsers, err := gorm.Open(postgres.Open(postgreSQL), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	UDB = connectUsers
	connectUsers.AutoMigrate(&Oauth.ModelOauth{})
}
