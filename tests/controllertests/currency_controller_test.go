package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pbutarbutar/game_currency/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateCurrency(t *testing.T) {

	err := refreshCurrencyTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"currency_from":1, "currency_to": 2, "rate":20, "author_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			name:         "Parulian1",
			author_id:    user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"currency_from":1, "currency_to": 2, "rate":20, "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/currencies", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateCurrency)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetCustomers(t *testing.T) {

	err := refreshCurrencyTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUCurrency()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/currencies", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetCustomer)
	handler.ServeHTTP(rr, req)

	var posts []models.Post
	err = json.Unmarshal([]byte(rr.Body.String()), &posts)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(posts), 2)
}
