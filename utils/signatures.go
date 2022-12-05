package utils

import (
	"encoding/json"
	"fmt"
	"rt-server/models"

	socketio "github.com/googollee/go-socket.io"
)

type RFunc func(c chan AppResponse)

func RFuncRun(cb RFunc) AppResponse {
	c := make(chan AppResponse)
	go cb(c)

	res := <-c
	return res
}

type ParseSUserFunc func(socketio.Conn, interface{}) AppResponse

func ParseSUserFuncRun(s socketio.Conn, d interface{}, cb ParseSUserFunc) AppResponse {
	keys, _ := RedisClient.Keys(RedisCtx, fmt.Sprintf("s_user_%s_*", s.ID())).Result()

	if len(keys) == 0 {
		s.Close()
		return ARFailure("Socket is not registered.")
	}

	key := keys[0]

	str, _ := RedisClient.Get(RedisCtx, key).Result()

	var sUser models.SUser

	json.Unmarshal([]byte(str), &sUser)

	s.SetContext(sUser)

	res := cb(s, d)

	return res
}