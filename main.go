package main

import (
	"final/configuration"
	"final/handler"
	"final/router"
	"log"
	"net/http"
)

func main() {
	// configuration
	newConfiguration := configuration.EnvConfiguration{}
	err := newConfiguration.Configure(".env")
	if err != nil {
		log.Fatal(err)
	}

	// handlers
	productHandler := handler.NewProductHandler()

	// router aka "server"
	routes := []router.HTTPRoute{
		{
			Method:  http.MethodGet,
			Pattern: "^/product/?$",
			Handler: productHandler.List,
		},
		{
			Method:  http.MethodGet,
			Pattern: "^/product/([^/]+)/?$",
			Handler: productHandler.Read,
		},
		{
			Method:  http.MethodPost,
			Pattern: "^/product/?$",
			Handler: productHandler.Create,
		},
		{
			Method:  http.MethodPatch,
			Pattern: "^/product/?$",
			Handler: productHandler.Update,
		},
	}
	router := router.NewHTTPRouter(newConfiguration, routes)

	err = router.Route()
	if err != nil {
		log.Fatal(err)
	}
}
