package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pbutarbutar/game_currency/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateCustomer(t *testing.T) {

	err := refreshUserAndCustomerTable()
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
			inputJSON:    `{"name":"The title", "author_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			name:         "Parulian1",
			author_id:    user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"Parulian1", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"name":"When no token is passed",  "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"name":"When incorrect token is passed",  "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"name": "",  "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			inputJSON:    `{"name": "This is an awesome title"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"name": "This is an awesome title", "content": "the content", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/customers", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateCustomer)

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

	err := refreshUserAndCustomerTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndCustomer()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/customers", nil)
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
func TestGetPostByID(t *testing.T) {

	err := refreshUserAndCustomerTable()
	if err != nil {
		log.Fatal(err)
	}
	cust, err := seedOneUserAndOneCust()
	if err != nil {
		log.Fatal(err)
	}
	custSample := []struct {
		id           string
		statusCode   int
		name         string
		author_id    uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(post.ID)),
			statusCode: 200,
			name:       cust.Name,
			author_id:  cust.AuthorID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range custSample {

		req, err := http.NewRequest("GET", "/customers", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetPost)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, cust.Name, responseMap["name"])
			assert.Equal(t, float64(post.AuthorID), responseMap["author_id"]) //the response author id is float64
		}
	}
}

func TestUpdateCustomer(t *testing.T) {

	var CustUserEmail, CustUserPassword string
	var AuthCustAuthorID uint32
	var AuthCustID uint64

	err := refreshUserAndCustomerTable()
	if err != nil {
		log.Fatal(err)
	}
	users, custs, err := seedUsersAndCustomer()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		CustUserEmail = user.Email
		CustUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(CustUserEmail, CustUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first post
	for _, cust := range custs {
		if cust.ID == 2 {
			continue
		}
		AuthCustID = cust.ID
		AuthCustAuthorID = cust.AuthorID
	}
	// fmt.Printf("this is the auth post: %v\n", AuthPostID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		name         string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"The updated post", "author_id": 1}`,
			statusCode:   200,
			name:         "The updated post",
			author_id:    AuthPostAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"This is still another title", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"This is still another title","author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Title 2" belongs to post 2, and title must be unique
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"Title 2", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"",  "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"Awesome title", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"This is another title"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"This is still another title", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/customers", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdatePost)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.title)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeletePost(t *testing.T) {

	var CustUserEmail, CustUserPassword string
	var AuthCustAuthorID uint32
	var AuthCustID uint64

	err := refreshUserAndCustomerTable()
	if err != nil {
		log.Fatal(err)
	}
	users, custs, err := seedUsersAndCustomer()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first post
	for _, cust := range custs {
		if cust.ID == 2 {
			continue
		}
		AuthCustID = cust.ID
		AuthCustAuthorID = cust.AuthorID
	}
	// fmt.Printf("this is the auth post: %v\n", AuthPostID)

	custSample := []struct {
		id           string
		updateJSON   string
		statusCode   int
		name         string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"The updated post", "author_id": 1}`,
			statusCode:   200,
			name:         "The updated post",
			author_id:    AuthPostAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"This is still another title", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"This is still another title","author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Title 2" belongs to post 2, and title must be unique
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"Title 2", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"",  "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"Awesome title", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"This is another title"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthCustID)),
			updateJSON:   `{"name":"This is still another title", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range custSample {

		req, _ := http.NewRequest("GET", "/customers", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeletePost)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
