package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/pbutarbutar/game_currency/api/controllers"
	"github.com/pbutarbutar/game_currency/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var customerInstance = models.Customer{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	log.Printf("Before calling m.Run() !!!")
	ret := m.Run()
	log.Printf("After calling m.Run() !!!")
	//os.Exit(m.Run())
	os.Exit(ret)
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "sqlite3" {
		//DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		testDbName := os.Getenv("TestDbName")
		server.DB, err = gorm.Open(TestDbDriver, testDbName)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
		server.DB.Exec("PRAGMA foreign_keys = ON")
	}

}

func refreshUserTable() error {
	server.DB.Exec("SET foreign_key_checks=0")
	err := server.DB.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	server.DB.Exec("SET foreign_key_checks=1")
	err = server.DB.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	log.Printf("refreshUserTable routine OK !!!")
	return nil
}

func seedOneUser() (models.User, error) {

	_ = refreshUserTable()

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err := server.DB.Debug().Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	log.Printf("seedOneUser routine OK !!!")
	return user, nil
}

func seedUsers() error {

	users := []models.User{
		models.User{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}

	for i := range users {
		err := server.DB.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}

	log.Printf("seedUsers routine OK !!!")
	return nil
}

func refreshUserAndPostTable() error {

	server.DB.Exec("SET foreign_key_checks=0")
	err := server.DB.Debug().DropTableIfExists(&models.Customer{}, &models.User{}).Error
	if err != nil {
		return err
	}
	server.DB.Exec("SET foreign_key_checks=1")
	err = server.DB.Debug().AutoMigrate(&models.User{}, &models.Customer{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	log.Printf("refreshUserAndPostTable routine OK !!!")
	return nil
}

func seedOneUserAndOnePost() (models.Post, error) {

	err := refreshUserAndPostTable()
	if err != nil {
		return models.Post{}, err
	}
	user := models.User{
		Nickname: "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err = server.DB.Debug().Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Customer{}, err
	}
	post := models.Customer{
		Name:     "This is the sam",
		AuthorID: user.ID,
	}
	err = server.DB.Debug().Model(&models.Customer{}).Create(&post).Error
	if err != nil {
		return models.Customer{}, err
	}

	log.Printf("seedOneUserAndOnePost routine OK !!!")
	return post, nil
}

func seedUsersAndCust() ([]models.User, []models.Customer, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Customer{}, err
	}
	var users = []models.User{
		models.User{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var custs = []models.Customer{
		models.Customer{
			name: "Name 1",
		},
		models.Customer{
			Name: "Name 2",
		},
	}

	for i := range users {
		err = server.DB.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		custs[i].AuthorID = users[i].ID

		err = server.DB.Debug().Model(&models.Customer{}).Create(&custs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	log.Printf("seedUsersAndPosts routine OK !!!")
	return users, pocustssts, nil
}
