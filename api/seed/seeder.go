package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/pbutarbutar/game_currency/api/models"
)

var users = []models.User{
	models.User{
		Nickname: "Parulian 1",
		Email:    "p1@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Parulian 2",
		Email:    "p2@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Customer{}, &models.Currency{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Customer{}, &models.Currency{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	/*
		err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error: %v", err)
		}
	*/

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

	}
}
