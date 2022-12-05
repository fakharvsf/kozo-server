package models

import (
	"fmt"

	"github.com/google/uuid"
)

type PersonalTaskReadOne struct {
	ID       uint `json:"id"`
}

func (personalTaskReadOne *PersonalTaskReadOne) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			return NullValidator("id", personalTaskReadOne.ID)
		},
	})

	return errors
}

type PersonalTaskCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (personalTaskCreate *PersonalTaskCreate) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			return NullValidator("title", personalTaskCreate.Title)
		},
		func() string {
			return NullValidator("description", personalTaskCreate.Description)
		},
	})

	return errors
}

type PersonalTaskUpdate struct {
	Status      string `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (personalTaskUpdate *PersonalTaskUpdate) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			if personalTaskUpdate.Status != "" {
				return PersonalTaskStatusValidator("status", personalTaskUpdate.Status)
			} else {
				return ""
			}
		},
	})

	return errors
}

type AuthLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (authLogin *AuthLogin) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			return EmailValidator("email", authLogin.Email)
		},
		func() string {
			return NullValidator("password", authLogin.Password)
		},
	})

	return errors
}

type AuthRegister struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (authRegister *AuthRegister) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			return NullValidator("first_name", authRegister.FirstName)
		},
		func() string {
			return NullValidator("last_name", authRegister.LastName)
		},
		func() string {
			return EmailValidator("email", authRegister.Email)
		},
		func() string {
			return LengthValidator("password", authRegister.Password, 8)
		},
	})

	return errors
}

type UserSendFriendRequest struct {
	UUID string `json:"uuid"`
}

func (userSendFriendRequest *UserSendFriendRequest) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			return NullValidator("uuid", userSendFriendRequest.UUID)
		},
		func() string {
			_, error := uuid.Parse(userSendFriendRequest.UUID)
			if error != nil {
				return fmt.Sprintf("%s is not valid.", "uuid")
			} else {
				return ""
			}
		},
	})

	return errors
}

type UserUpdateFriendRequest struct {
	ID uint `json:"id"`
	Status int `json:"status"`
}

func (userUpdateFriendRequest *UserUpdateFriendRequest) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			return IDValidator("id", userUpdateFriendRequest.ID)
		},
		func() string {
			return NullValidator("status", userUpdateFriendRequest.Status)
		},
	})

	return errors
}

type UserSearchRequest struct {
	UUID string `json:"uuid"`
}

func (userSearchRequest *UserSearchRequest) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			return NullValidator("uuid", userSearchRequest.UUID)
		},
	})

	fmt.Println("errors", errors)

	return errors
}

type AssignPersonalTask struct {
	PersonalTaskID uint `json:"personal_task_id"`
	AssignedToID uint `json:"assigned_to_id"`
}

func (assignPersonalTask *AssignPersonalTask) Validate() []string {
	errors := Validator([]ValidatorClosure{
		func() string {
			return IDValidator("personal_task_id", assignPersonalTask.PersonalTaskID)
		},
		func() string {
			return IDValidator("assigned_to_id", assignPersonalTask.AssignedToID)
		},
	})

	return errors
}

type PersonalAndAssignedTasks struct {
	PersonalTasks []PersonalTaskJSON `json:"personal_tasks"`
	AssignedTasks []AssignedTaskJSON `json:"assigned_tasks"`
}
