package auth

import (
	"kozo/models"
	"kozo/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(authLogin models.AuthLogin, c chan utils.AppResponse) {
	valError := authLogin.Validate()

	if len(valError) > 0 {
		c <- utils.ARValFail("", valError)
		return
	}

	var user models.User

	dbRes := utils.DB.Where("email = ?", authLogin.Email).First(&user)

	if dbRes.Error != nil {
		c <- utils.ARFailure("User not found.")
		return
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authLogin.Password))

	if passwordError != nil {
		c <- utils.ARFailure("Password is not correct.")
		return
	}

	token := utils.GenerateJwtToken(user.ID, utils.UserRolesEnums["user"])
	tokenStr, _ := token.SignedString([]byte(utils.JwtSigningKey))

	userJSON := user.ToJSON()
	userJSON.Token = tokenStr

	c <- utils.ARSuccess(userJSON)
}

func Register(authRegister models.AuthRegister, c chan utils.AppResponse){
	valError := authRegister.Validate()

	if len(valError) > 0 {
		c <- utils.ARValFail("", valError)
		return
	}

	var existingUser models.User

	existingUserError := utils.DB.Where("email = ?", authRegister.Email).First(&existingUser)

	if existingUserError.Error == nil {
		c <- utils.ARFailure("User with this email already exists.")
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(authRegister.Password), 10)
	uuid := uuid.New().String()

	user := models.User{
		FirstName: authRegister.FirstName,
		LastName: authRegister.LastName,
		Email: authRegister.Email,
		Password: string(password),
		IsActive: true,
		UUID: uuid,
	}

	utils.DB.Create(&user)

	token := utils.GenerateJwtToken(user.ID, utils.UserRolesEnums["user"])
	tokenStr, _ := token.SignedString([]byte(utils.JwtSigningKey))

	userJSON := user.ToJSON()
	userJSON.Token = tokenStr

	c <- utils.ARSuccess(userJSON)
}

func Sync(ID uint, c chan utils.AppResponse) {
	var user models.User

	dbRes := utils.DB.Where("id = ?", ID).First(&user)

	if dbRes.Error != nil {
		c <- utils.ARFailure("User not found.")
		return
	}

	token := utils.GenerateJwtToken(user.ID, utils.UserRolesEnums["user"])
	tokenStr, _ := token.SignedString([]byte(utils.JwtSigningKey))

	userJSON := user.ToJSON()
	userJSON.Token = tokenStr

	c <- utils.ARSuccess(userJSON)
}