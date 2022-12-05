package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	FirstName string    `gorm:"column:first_name;not null"`
	LastName  string	`gorm:"column:last_name;not null"`
	Email     string	`gorm:"column:email;not null"`
	Password  string	`gorm:"column:password;not null"`
	IsActive  bool		`gorm:"column:is_active;not null"`
	UUID      string	`gorm:"column:uuid;not null"`
	CreatedAt time.Time	`gorm:"column:created_at;not null"`
	UpdatedAt time.Time	`gorm:"column:updated_at;not null"`
}

type UserJSON struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	UUID	  string	`json:"uuid"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (user *User) ToJSON() UserJSON {
	return UserJSON{
		ID: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		UUID: user.UUID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}	
}

type PersonalTask struct {
	ID          uint      `gorm:"column:id;primaryKey"`
	Title       string	  `gorm:"column:title;not null"`
	Description string    `gorm:"column:description;not null"`
	Status      string	  `gorm:"column:status;not null;default:created"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
	UserID      uint	  `gorm:"index;not null"`
	User        User      `gorm:"foreignKey:UserID"`
}

type PersonalTaskJSON struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        UserJSON  `json:"user"`
}

func (personalTask *PersonalTask) ToJSON() PersonalTaskJSON {
	return PersonalTaskJSON{
		ID: personalTask.ID,
		Title: personalTask.Title,
		Description: personalTask.Description,
		Status: personalTask.Status,
		CreatedAt: personalTask.CreatedAt,
		UpdatedAt: personalTask.UpdatedAt,
		User: personalTask.User.ToJSON(),
	}
}

type FriendRequest struct {
	ID          uint      `gorm:"column:id;primaryKey"`
	SenderID    uint	  `gorm:"column:sender_id;index;not null"`
	Sender		User	  `gorm:"foreignKey:SenderID"`
	ReceiverID  uint	  `gorm:"column:receiver_id;index;not null"`
	Receiver	User	  `gorm:"foreignKey:ReceiverID"`
	Status      string	  `gorm:"column:status;not null;default:sent"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

type FriendRequestJSON struct {
	ID          uint      `json:"id"`
	Sender      UserJSON  `json:"sender"`
	Receiver    UserJSON  `json:"receiver"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (friendRequest *FriendRequest) ToJSON() FriendRequestJSON {
	return FriendRequestJSON{
		ID: friendRequest.ID,
		Sender: friendRequest.Sender.ToJSON(),
		Receiver: friendRequest.Receiver.ToJSON(),
		Status: friendRequest.Status,
		CreatedAt: friendRequest.CreatedAt,
		UpdatedAt: friendRequest.UpdatedAt,
	}
}

type Friend struct {
	ID          uint      `gorm:"column:id;primaryKey"`
	SenderID    uint	  `gorm:"column:sender_id;index;not null"`
	Sender		User	  `gorm:"foreignKey:SenderID"`
	ReceiverID  uint	  `gorm:"column:receiver_id;index;not null"`
	Receiver	User	  `gorm:"foreignKey:ReceiverID"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

type FriendJSON struct {
	ID          uint      `json:"id"`
	Sender      UserJSON  `json:"sender"`
	Receiver    UserJSON  `json:"receiver"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (friend *Friend) ToJSON() FriendJSON {
	return FriendJSON{
		ID: friend.ID,
		Sender: friend.Sender.ToJSON(),
		Receiver: friend.Receiver.ToJSON(),
		CreatedAt: friend.CreatedAt,
		UpdatedAt: friend.UpdatedAt,
	}
}

type AssignedTask struct {
	ID             uint         `gorm:"column:id;primaryKey"`
	OwnerID        uint	        `gorm:"column:owner_id;index;not null"`
	Owner		   User	        `gorm:"foreignKey:OwnerID"`
	AsssignedToID  uint	        `gorm:"column:assigned_to_id;index;not null"`
	AssignedTo	   User	        `gorm:"foreignKey:AsssignedToID"`
	PersonalTaskID uint	        `gorm:"column:personal_task_id;index;not null"`
	PersonalTask   PersonalTask `gorm:"foreignKey:PersonalTaskID"`
	CreatedAt      time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time    `gorm:"column:updated_at;not null"`
}

type AssignedTaskJSON struct {
	ID           uint             `json:"id"`
	Owner        UserJSON         `json:"owner"`
	AssignedTo   UserJSON         `json:"assigned_to"`
	PersonalTask PersonalTaskJSON `json:"personal_task"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func (assignedTask *AssignedTask) ToJSON() AssignedTaskJSON {
	return AssignedTaskJSON{
		ID: assignedTask.ID,
		Owner: assignedTask.Owner.ToJSON(),
		AssignedTo: assignedTask.AssignedTo.ToJSON(),
		PersonalTask: assignedTask.PersonalTask.ToJSON(),
		CreatedAt: assignedTask.CreatedAt,
		UpdatedAt: assignedTask.UpdatedAt,
	}
}

