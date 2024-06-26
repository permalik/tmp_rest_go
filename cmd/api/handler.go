package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/permalik/temp_rest_go/internal/data"
	"github.com/permalik/temp_rest_go/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.err_res_notf(w, r)
		return
	}
	fmt.Fprintln(w, "home...")
}

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	data := wrap_json{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			// TODO: convert version to env var
			"version": "0.1.0",
		},
	}
	err := app.w_json(w, http.StatusOK, data, nil, true)
	if err != nil {
		app.err_res_int_srv(w, r, err)
	}
}

func (app *application) create_item(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string      `json:"name"`
		Year     int32       `json:"year"`
		Quantity int32       `json:"quantity"`
		Pounds   data.Pounds `json:"pounds"`
		Types    []string    `json:"types"`
	}
	err := app.r_json(w, r, &input)
	if err != nil {
		app.err_res_bad_req(w, r, err)
	}

	i := &data.Item{
		Name:     input.Name,
		Year:     input.Year,
		Quantity: input.Quantity,
		Pounds:   input.Pounds,
		Types:    input.Types,
	}

	v := validator.New()

	if data.ValidateItem(v, i); !v.Valid() {
		app.err_res_fail_val(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) read_items(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "read items...")
}

func (app *application) read_item(w http.ResponseWriter, r *http.Request) {
	id, err := app.parse_key(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	data := data.Item{
		ID:        id,
		Name:      "test_item",
		Year:      1986,
		Quantity:  12,
		Pounds:    15,
		Types:     []string{"primary", "secondary", "tertiary"},
		CreatedAt: time.Now(),
	}

	err = app.w_json(w, http.StatusOK, wrap_json{"item": data}, nil, true)
	if err != nil {
		app.err_res_int_srv(w, r, err)
	}
}
