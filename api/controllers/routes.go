package controllers

import "github.com/pbutarbutar/game_currency/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Customer Routes
	s.Router.HandleFunc("/customers", middlewares.SetMiddlewareJSON(s.CreateCustomer)).Methods("POST")
	s.Router.HandleFunc("/customers", middlewares.SetMiddlewareJSON(s.GetCustomers)).Methods("GET")
	s.Router.HandleFunc("/customers/{id}", middlewares.SetMiddlewareJSON(s.GetCustomer)).Methods("GET")
	s.Router.HandleFunc("/customers/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCustomer))).Methods("PUT")
	s.Router.HandleFunc("/customers/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCustomer)).Methods("DELETE")

	//Currency Routes
	s.Router.HandleFunc("/calculatecurrency", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CalculateCurrency))).Methods("POST")
	s.Router.HandleFunc("/currencies", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateCurrency))).Methods("POST")
	s.Router.HandleFunc("/currencies", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCurrencies))).Methods("GET")
	s.Router.HandleFunc("/currencies/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCurrency))).Methods("GET")
	s.Router.HandleFunc("/currencies/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCurrency))).Methods("PUT")
	s.Router.HandleFunc("/currencies/{id}", middlewares.SetMiddlewareAuthentication(middlewares.SetMiddlewareAuthentication(s.DeleteCurrency))).Methods("DELETE")

}
