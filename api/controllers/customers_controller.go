package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pbutarbutar/game_currency/api/auth"
	"github.com/pbutarbutar/game_currency/api/models"
	"github.com/pbutarbutar/game_currency/api/responses"
	"github.com/pbutarbutar/game_currency/api/utils/formaterror"
)

func (server *Server) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	customer := models.Customer{}
	err = json.Unmarshal(body, &customer)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	customer.Prepare()
	err = customer.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != customer.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	xCheck := customer.IsCheckCustExist(server.DB, customer.Name)
	fmt.Println(xCheck)
	if xCheck > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Name is exist"))
		return
	}

	customerCreated, err := customer.SaveCustomer(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, customerCreated.ID))
	responses.JSON(w, http.StatusCreated, customerCreated)
}

func (server *Server) GetCustomers(w http.ResponseWriter, r *http.Request) {

	customer := models.Customer{}

	customers, err := customer.FindAllCustomer(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, customers)
}

func (server *Server) GetCustomer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	fmt.Println(pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	customer := models.Customer{}

	customerReceived, err := customer.FindCustomerByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, customerReceived)
}

func (server *Server) UpdateCustomer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	cust := models.Customer{}
	err = server.DB.Debug().Model(models.Customer{}).Where("id = ?", pid).Take(&cust).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	// If a user attempt to update a post not belonging to him
	if uid != cust.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	custUpdate := models.Customer{}
	err = json.Unmarshal(body, &custUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != custUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	custUpdate.Prepare()
	err = custUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	custUpdate.ID = cust.ID //this is important to tell the model the post id to update, the other update field are set above

	custUpdated, err := custUpdate.UpdateAXCustomer(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, custUpdated)
}

func (server *Server) DeleteCustomer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	customer := models.Customer{}
	err = server.DB.Debug().Model(models.Customer{}).Where("id = ?", cid).Take(&customer).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != customer.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = customer.DeleteACustomer(server.DB, cid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", cid))
	responses.JSON(w, http.StatusNoContent, "")
}
