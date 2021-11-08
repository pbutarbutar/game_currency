package modeltests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pbutarbutar/game_currency/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	_, _, err = seedUsersAndCust()
	if err != nil {
		log.Fatalf("Error seeding user and post table %v\n", err)
	}
	posts, err := customerInstance.FindAllCustomer()(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*posts), 2)
}

func TestSavePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newPost := models.Customer{
		ID:       1,
		Name:     "This is the Name",
		AuthorID: user.ID,
	}
	savedPost, err := newPost.SaveCustomer(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newPost.ID, savedPost.ID)
	assert.Equal(t, newPost.Name, savedPost.Name)

}

func TestGetPostByID(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding User and post table")
	}
	foundPost, err := customerInstance.FindCustomerByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundPost.ID, post.ID)
	assert.Equal(t, foundPost.Name, post.Name)
}
