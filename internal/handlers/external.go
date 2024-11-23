package handlers

import "net/http"

// ExampleHandler godoc
// @Summary Example endpoint
// @Description Responds with a simple message
// @Tags example
// @Success 200 {string} string "OK"
// @Router /example [get]
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, Swagger!"))
}