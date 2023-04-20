package socket

import (
	"encoding/json"
	"fmt"
	"kozo/controllers/tasks"
	"kozo/models"
	"kozo/utils"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func Connect(s socketio.Conn) error {
	s.SetContext("")
	fmt.Println("Connected:", s.ID())
	return nil
}

func DisConnect(s socketio.Conn, data string) {
	keys, _ := utils.RedisClient.Keys(utils.RedisCtx, fmt.Sprintf("s_user_%s_*", s.ID())).Result()
	if len(keys) > 0 {
		key := keys[0]
		utils.RedisClient.Del(utils.RedisCtx, key)
		fmt.Println("Redis Removed:", key)
	}
	fmt.Println("Disconnected:", s.ID())
}

func Register(s socketio.Conn, d interface{}) utils.AppResponse {
	dataStr, _ := json.Marshal(d)

	var srRegister models.SRRegister

	dataErr := json.Unmarshal(dataStr, &srRegister)
	if dataErr != nil {
		s.Close()
		return utils.ARFailure("Could not parse data");
	}

	valError := srRegister.Validate()
	if len(valError) > 0 {
		s.Close()
		return utils.ARValFail("", valError);
	}

	tokenFullStrs := strings.Split(srRegister.Token, " ")

	if len(tokenFullStrs) <= 1 {
		s.Close()
		return utils.ARFailure("Bearer token not found.");
	}

	token := tokenFullStrs[1]

	jwtClaims, jwtClaimsError := utils.ParseJwtToken(token)

	if jwtClaimsError != nil {
		s.Close()
		return utils.ARFailure("Could not parse token.");
	}

	ID := jwtClaims.ID

	var user models.User

	utils.DB.Where("id = ?", ID).First(&user)

	sUser := models.SUser{
		ID: s.ID(),
		User: user.ToJSON(),
	}

	existingKeys, _:= utils.RedisClient.Keys(utils.RedisCtx, fmt.Sprintf("s_user_*_%d", sUser.User.ID)).Result()

	if len(existingKeys) > 0 {
		s.Close()
		return utils.ARFailure("This user already has a connection.");
	}

	userStr, _ := json.Marshal(sUser)
	utils.RedisClient.Set(utils.RedisCtx, fmt.Sprintf("s_user_%s_%d", sUser.ID, sUser.User.ID), string(userStr), 0)

	return utils.ARSuccess("Successfully registered!")
}

func PersonalTasksCreate(s socketio.Conn, d interface{}) utils.AppResponse {
	res := utils.RFuncRun(
		func (c chan utils.AppResponse){
			sUser := s.Context().(models.SUser)

			personalTaskCreateStr, _ := json.Marshal(d)
			var personalTaskCreate models.PersonalTaskCreate
			err := json.Unmarshal(personalTaskCreateStr, &personalTaskCreate)
			if err != nil {
				c <- utils.ARFailure("Incorrect data received.")
				return
			}
			valError := personalTaskCreate.Validate()
			if len(valError) > 0 {
				c <- utils.ARValFail("", valError)
				return
			}

			tasks.Create(sUser.User.ID, personalTaskCreate, c)
		},
	)

	return res
}

func PersonalTasksReadOne(s socketio.Conn, d interface{}) utils.AppResponse {
	res := utils.RFuncRun(
		func (c chan utils.AppResponse){
			sUser := s.Context().(models.SUser)

			personalTaskReadOneStr, _ := json.Marshal(d)
			var personalTaskReadOne models.PersonalTaskReadOne
			err := json.Unmarshal(personalTaskReadOneStr, &personalTaskReadOne)
			if err != nil {
				c <- utils.ARFailure("Incorrect data received.")
				return
			}
			valError := personalTaskReadOne.Validate()
			if len(valError) > 0 {
				c <- utils.ARValFail("", valError)
				return
			}

			tasks.ReadOne(sUser.User.ID, personalTaskReadOne.ID, c)
		},
	)

	return res
}

func PersonalTasksReadAll(s socketio.Conn, d interface{}) utils.AppResponse {
	res := utils.RFuncRun(
		func (c chan utils.AppResponse){
			sUser := s.Context().(models.SUser)
			tasks.ReadAll(sUser.User.ID, c)
		},
	)

	return res
}

func AssignedTasksReadAll(s socketio.Conn, d interface{}) utils.AppResponse {
	res := utils.RFuncRun(
		func (c chan utils.AppResponse){
			sUser := s.Context().(models.SUser)
			tasks.ReadAssignedTasks(sUser.User.ID, c)
		},
	)

	return res
}

func TasksReadAll(s socketio.Conn, d interface{}) utils.AppResponse {
	res := utils.RFuncRun(
		func (c chan utils.AppResponse){
			sUser := s.Context().(models.SUser)
			tasks.ReadPersonalAndAssigned(sUser.User.ID, c)
		},
	)

	return res
}
