package main

import (
	"fmt"
	"kozo/utils"
	"kozo/views/socket"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
)

func MainSocketServer(env string) {
	// Load Env
	err := godotenv.Load(env)

	if err != nil {
		panic(err)
	}
	
	// DB Connection
	dbError := utils.DBConnect()
	if dbError != nil {
		panic(dbError)
	}

	// Run migrations
	utils.DBMigrate(false)

	// Redis Connection
	redisError := utils.RedisConnect()
	if redisError != nil {
		panic(redisError)
	}

	server := socketio.NewServer(nil)

	_, redisSocketAdapterErr := server.Adapter(&socketio.RedisAdapterOptions{
		Addr:   "127.0.0.1:6379",
		Host:   "127.0.0.1",
		Port:   "6379",
		Prefix: "socket.io",
	})
	if redisSocketAdapterErr != nil {
		panic(redisSocketAdapterErr)
	}

	server.OnConnect("/", socket.Connect)
	server.OnDisconnect("/", socket.DisConnect)
	server.OnEvent("/", "register", socket.Register)

	server.OnEvent("/", "tasks:create", func (s socketio.Conn, d interface{}) utils.AppResponse {
		res := utils.ParseSUserFuncRun(s, d, socket.PersonalTasksCreate)
		return res
	})

	server.OnEvent("/", "tasks:personalReadAll", func (s socketio.Conn, d interface{}) utils.AppResponse {
		res := utils.ParseSUserFuncRun(s, d, socket.PersonalTasksReadAll)
		return res
	})
	
	server.OnEvent("/", "tasks:personalReadOne", func (s socketio.Conn, d interface{}) utils.AppResponse {
		res := utils.ParseSUserFuncRun(s, d, socket.PersonalTasksReadOne)
		return res
	})

	server.OnEvent("/", "tasks:assignedReadAll", func (s socketio.Conn, d interface{}) utils.AppResponse {
		res := utils.ParseSUserFuncRun(s, d, socket.AssignedTasksReadAll)
		return res
	})

	server.OnEvent("/", "tasks:readAll", func (s socketio.Conn, d interface{}) utils.AppResponse {
		res := utils.ParseSUserFuncRun(s, d, socket.TasksReadAll)
		return res
	})

	go server.Serve()
	defer server.Close()

	serverIP := os.Getenv("serverIP");
	serverPort := os.Getenv("serverPort");
	fmt.Println("Socket server is up and running on Port:", serverPort)
	http.ListenAndServe(fmt.Sprintf("%s:%s", serverIP, serverPort), server)
}