package controllers

import (
	"net/http"

	"github.com/pbutarbutar/game_currency/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "API GOLANG -  PARULIAN BUTAR BUTAR")

}
