package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pbutarbutar/game_currency/api/models"
	"github.com/pbutarbutar/game_currency/api/responses"
	"github.com/pbutarbutar/game_currency/api/utils/formaterror"
)

func (server *Server) CalculateCurrency(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	calccurrency := models.RequestCalCurrency{}
	err = json.Unmarshal(body, &calccurrency)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var resCurr models.ResponseCalCurrency
	currency := models.Currency{}

	currencyReceived, err := currency.FindCurrencyByCurrencyID(server.DB, calccurrency.CurrencyFrom, calccurrency.CurrencyTo)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	resCurr.CurrencyFrom = currencyReceived.CurrencyFrom
	resCurr.CurrencyTo = currencyReceived.CurrencyTo
	resCurr.Amount = calccurrency.Amount
	resCurr.Result = calccurrency.Amount / currencyReceived.Rate
	responses.JSON(w, http.StatusOK, resCurr)

}

func (server *Server) CreateCurrency(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	currency := models.Currency{}
	err = json.Unmarshal(body, &currency)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	currency.Prepare()
	err = currency.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	xCheck := currency.IsCheckExist(server.DB, currency.CurrencyFrom, currency.CurrencyTo)

	if xCheck > 0 {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Currency Is exist"))
		return
	}

	currencyCreated, err := currency.SaveCurrency(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, currencyCreated.ID))
	responses.JSON(w, http.StatusCreated, currencyCreated)
}

func (server *Server) GetCurrencies(w http.ResponseWriter, r *http.Request) {

	currency := models.Currency{}

	currencies, err := currency.FindAllCurrency(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, currencies)
}

func (server *Server) GetCurrency(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	fmt.Println(pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	currency := models.Currency{}

	currencyReceived, err := currency.FindCurrencyByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, currencyReceived)
}

func (server *Server) UpdateCurrency(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if the post exist
	cust := models.Currency{}
	err = server.DB.Debug().Model(models.Currency{}).Where("id = ?", pid).Take(&cust).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	custUpdate := models.Currency{}
	err = json.Unmarshal(body, &custUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	custUpdate.Prepare()
	err = custUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	custUpdate.ID = cust.ID //this is important to tell the model the post id to update, the other update field are set above

	currUpdated, err := custUpdate.UpdateACurrency(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, currUpdated)
}

func (server *Server) DeleteCurrency(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if the post exist
	curr := models.Currency{}
	err = server.DB.Debug().Model(models.Currency{}).Where("id = ?", cid).Take(&curr).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = curr.DeleteACurrency(server.DB, cid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", cid))
	responses.JSON(w, http.StatusNoContent, "")
}
