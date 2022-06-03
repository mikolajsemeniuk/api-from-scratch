package handler

import (
	"encoding/json"
	"final/router"
	"final/validator"
	"net/http"
)

type ProductHandler struct{}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Product list"))
}

func (h *ProductHandler) Read(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(router.ContextParamsKey).([]string)
	w.Write([]byte("Product read: " + params[0]))
}

func (*ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string  `json:"name"        re:"^.{4,8}$"`
		Description string  `json:"description" re:"^.{5,25}$"`
		Price       float32 `json:"price"       range:",25"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Validate(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (*ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        *string  `json:"name,omitempty" re:"^.{4,8}$"`
		Description *string  `json:"description,omitempty" re:"^.{5,25}$"`
		Price       *float32 `json:"price,omitempty" range:",300"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Validate(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (*ProductHandler) Remove(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func NewProductHandler() Handler {
	return &ProductHandler{}
}
