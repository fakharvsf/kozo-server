package rest

import (
	"encoding/json"
	"net/http"
	"rt-server/controllers/users"
	"rt-server/models"
	"rt-server/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func Users(r chi.Router) {
	r.Use(utils.ParseToken)

	r.Get("/friend_requests", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				ID := r.Context().Value("ID").(uint)
				users.GetFriendRequests(ID, c)
			},
		)
		render.JSON(w, r, res)
	})
	r.Post("/friend_requests", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				ID := r.Context().Value("ID").(uint)
				var userSendFriendRequest models.UserSendFriendRequest
				json.NewDecoder(r.Body).Decode(&userSendFriendRequest)

				users.SendFriendRequest(ID, userSendFriendRequest, c)
			},
		)
		render.JSON(w, r, res)
	})
	r.Put("/friend_requests", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				ID := r.Context().Value("ID").(uint)
				var userUpdateFriendRequest models.UserUpdateFriendRequest
				json.NewDecoder(r.Body).Decode(&userUpdateFriendRequest)

				users.UpdateFriendRequest(ID, userUpdateFriendRequest, c)
			},
		)
		render.JSON(w, r, res)
	})
	r.Get("/friends", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				ID := r.Context().Value("ID").(uint)

				users.GetFriends(ID, c)
			},
		)
		render.JSON(w, r, res)
	})
	r.Post("/search", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				var userSearchRequest models.UserSearchRequest
				json.NewDecoder(r.Body).Decode(&userSearchRequest)

				users.Search(userSearchRequest, c)
			},
		)
		render.JSON(w, r, res)
	})
}