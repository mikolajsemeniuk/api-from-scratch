package router

import (
	"context"
	"final/configuration"
	"net/http"
	"regexp"
)

type ParamsKey string

const ContextParamsKey ParamsKey = "params"

type HTTPRoute struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

type HTTPRouter struct {
	configuration configuration.EnvConfiguration
	routes        []HTTPRoute
}

func (r *HTTPRouter) Route() error {
	http.HandleFunc("/", r.Handler)
	return http.ListenAndServe(r.configuration.Listen, nil)
}

func (r *HTTPRouter) Handler(writer http.ResponseWriter, request *http.Request) {
	url := request.URL.Path

	for _, route := range r.routes {
		re := regexp.MustCompile(route.Pattern)
		matches := re.FindAllStringSubmatch(url, -1)

		if len(matches) != 1 {
			continue
		}

		if request.Method == route.Method {
			params := matches[0][1:]
			newContext := context.WithValue(request.Context(), ContextParamsKey, params)
			route.Handler(writer, request.WithContext(newContext))
			return
		}
	}

	writer.WriteHeader(http.StatusNotFound)
}

func NewHTTPRouter(configuration configuration.EnvConfiguration, routes []HTTPRoute) Router {
	return &HTTPRouter{
		configuration: configuration,
		routes:        routes,
	}
}
