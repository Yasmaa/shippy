package service

import "net/http"

type Route struct {

	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{

	Route{
		"GetAccount",                                     // Name
		"GET",                                            // HTTP method
		"/accounts/{accountId}",                          // Route pattern
		GetAccount,
	},
	Route{
        "HealthCheck",
        "GET",
        "/accounts/health",
        HealthCheck,
},
}