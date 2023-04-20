package rest

import (
	"encoding/json"
	"kozo/controllers/tasks"
	"kozo/models"
	"kozo/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func Tasks(r chi.Router) {
	r.Use(utils.ParseToken)

	r.Get("/", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				ID := r.Context().Value("ID").(uint)

				tasks.ReadAll(ID, c)
			},
		)

		render.JSON(w, r, res)
	})
	r.Post("/", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				ID := r.Context().Value("ID").(uint)
				var personalTaskCreate models.PersonalTaskCreate
				json.NewDecoder(r.Body).Decode(&personalTaskCreate)

				tasks.Create(ID, personalTaskCreate, c)
			},
		)
		render.JSON(w, r, res)
	})

	r.Get("/all", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				ID := r.Context().Value("ID").(uint)

				tasks.ReadPersonalAndAssigned(ID, c)
			},
		)

		render.JSON(w, r, res)
	})
	

	r.Route("/{id}", func(r chi.Router){
		r.Get("/", func (w http.ResponseWriter, r *http.Request){
			res := utils.RFuncRun(
				func (c chan utils.AppResponse){
					ID := r.Context().Value("ID").(uint)
					personalTaskID, errorID := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
					if errorID != nil {
						c <- utils.ARFailure(errorID.Error())
						return
					}

					tasks.ReadOne(ID, uint(personalTaskID), c)
				},
			)
			render.JSON(w, r, res)
		})
		r.Put("/", func (w http.ResponseWriter, r *http.Request){
			res := utils.RFuncRun(
				func (c chan utils.AppResponse){
					ID := r.Context().Value("ID").(uint)
					personalTaskID, errorID := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
					if errorID != nil {
						c <- utils.ARFailure(errorID.Error())
						return
					}
					var personalTaskUpdate models.PersonalTaskUpdate
					json.NewDecoder(r.Body).Decode(&personalTaskUpdate)

					tasks.Update(ID, personalTaskID, personalTaskUpdate, c)
				},
			)
			render.JSON(w, r, res)
		})
		r.Delete("/", func (w http.ResponseWriter, r *http.Request){
			res := utils.RFuncRun(
				func (c chan utils.AppResponse){
					ID, errorID := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
					if errorID != nil {
						c <- utils.ARFailure(errorID.Error())
						return
					}

					tasks.Delete(ID, c)
				},
			)
			render.JSON(w, r, res)
		})
	})

	// Assignees
	r.Get("/read_assigned", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				ID := r.Context().Value("ID").(uint)
				tasks.ReadAssignedTasks(ID, c)
			},
		)
		render.JSON(w, r, res)
	})
	r.Post("/assign", func (w http.ResponseWriter, r *http.Request){
		res := utils.RFuncRun(
			func (c chan utils.AppResponse){
				ID := r.Context().Value("ID").(uint)
				var assignPersonalTask models.AssignPersonalTask
				json.NewDecoder(r.Body).Decode(&assignPersonalTask)

				tasks.Assign(ID, assignPersonalTask, c)
			},
		)
		render.JSON(w, r, res)
	})
}