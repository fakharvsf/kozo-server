package rest

import (
	"encoding/json"
	"net/http"
	"rt-server/controllers/auth"
	"rt-server/models"
	"rt-server/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func Auth(r chi.Router) {
	r.Post("/login", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				var authLogin models.AuthLogin
				json.NewDecoder(r.Body).Decode(&authLogin)

				auth.Login(authLogin, c)
			},
		)
		render.JSON(w, r, res)
	})
	r.Post("/register", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				var authRegister models.AuthRegister
				json.NewDecoder(r.Body).Decode(&authRegister)

				auth.Register(authRegister, c)
			},
		)
		render.JSON(w, r, res)
	})
	r.Group(func(r chi.Router) {
		r.Use(utils.ParseToken)
		r.Get("/sync", func (w http.ResponseWriter, r *http.Request){
			res := utils.RFuncRun(
				func (c chan utils.AppResponse){
					ID := r.Context().Value("ID").(uint)

					auth.Sync(ID, c)
				},
			)
			render.JSON(w, r, res)
		})
	})
}