package users

import (
	"kozo/models"
	"kozo/utils"
)

func GetFriendRequests(ID uint, c chan utils.AppResponse) {
	friendRequests := []models.FriendRequest{}

	utils.DB.Preload("Sender").Preload("Receiver").Where("receiver_id = ? AND status = ?", ID, models.FriendRequestsStatusesEnums["sent"]).Find(&friendRequests)

	friendRequestsParsed := []models.FriendRequestJSON{}

	for i := 0; i < len(friendRequests); i++ {
		friendRequest := friendRequests[i]
		friendRequestsParsed = append(friendRequestsParsed, friendRequest.ToJSON())
	} 

	c <- utils.ARSuccess(friendRequestsParsed)
}

func SendFriendRequest(ID uint, userSendFriendRequest models.UserSendFriendRequest, c chan utils.AppResponse) {
	valError := userSendFriendRequest.Validate()

	if len(valError) > 0 {
		c <- utils.ARValFail("", valError)
		return
	}

	var sender models.User

	utils.DB.First(&sender, ID)

	var receiver models.User

	dbRes := utils.DB.Where("uuid = ?", userSendFriendRequest.UUID).First(&receiver)

	if dbRes.Error != nil {
		c <- utils.ARFailure("User not found.")
		return
	}

	if userSendFriendRequest.UUID == sender.UUID {
		c <- utils.ARFailure("Users can't send friend requests to themselves.")
		return
	}

	existingFriendRes := utils.DB.Where("(sender_id = ? AND receiver_id = ?) OR (receiver_id = ? AND sender_id = ?)", sender.ID, receiver.ID, sender.ID, receiver.ID).First(&models.Friend{})

	if existingFriendRes.Error == nil {
		c <- utils.ARFailure("You are already friend with this user.")
		return
	}

	existingFriendRequestRes := utils.DB.Where("(sender_id = ? AND receiver_id = ?) OR (receiver_id = ? AND sender_id = ?)", sender.ID, receiver.ID, sender.ID, receiver.ID).First(&models.FriendRequest{})

	if existingFriendRequestRes.Error == nil {
		c <- utils.ARFailure("Friend request already exists.")
		return
	}

	friendRequest := models.FriendRequest{
		SenderID: ID,
		ReceiverID: receiver.ID,
		Status: models.FriendRequestsStatusesEnums["sent"],
	}

	utils.DB.Create(&friendRequest)

	friendRequest.Sender = sender
	friendRequest.Receiver = receiver

	c <- utils.ARSuccess(friendRequest.ToJSON())
}

func UpdateFriendRequest(ID uint, userUpdateFriendRequest models.UserUpdateFriendRequest, c chan utils.AppResponse) {
	valError := userUpdateFriendRequest.Validate()

	if len(valError) > 0 {
		c <- utils.ARValFail("", valError)
		return
	}

	var friendRequest models.FriendRequest

	dbRes := utils.DB.First(&friendRequest, userUpdateFriendRequest.ID)

	if dbRes.Error != nil {
		c <- utils.ARFailure("Friend request not found.")
		return
	}

	if friendRequest.ReceiverID != ID {
		c <- utils.ARFailure("You are not authorized to perform this action.")
		return
	}

	if friendRequest.Status != models.FriendRequestsStatusesEnums["sent"] {
		c <- utils.ARFailure("This friend request has already been upated.")
		return
	}

	var status string
	if userUpdateFriendRequest.Status == 0 {
		status = models.FriendRequestsStatusesEnums["rejected"]
	} else {
		status = models.FriendRequestsStatusesEnums["accepted"]
	}

	friendRequest.Status = status

	utils.DB.Save(&friendRequest)

	if userUpdateFriendRequest.Status == 0 {
		c <- utils.ARSuccess("Friend request rejected.")
	} else {
		friend := models.Friend{
			SenderID: friendRequest.SenderID,
			ReceiverID: friendRequest.ReceiverID,
		}
	
		utils.DB.Create(&friend)

		c <- utils.ARSuccess("Friend request accepted.")
	}
}

func GetFriends(ID uint, c chan utils.AppResponse) {
	friends := []models.Friend{}

	utils.DB.Preload("Sender").Preload("Receiver").Where("receiver_id = ? OR sender_id = ?", ID, ID).Find(&friends)

	friendsParsed := []models.UserJSON{}

	for i := 0; i < len(friends); i++ {
		friend := friends[i]
		if friend.Sender.ID == ID {
			friendsParsed = append(friendsParsed, friend.Receiver.ToJSON())
		} else {
			friendsParsed = append(friendsParsed, friend.Sender.ToJSON())
		}
	} 

	c <- utils.ARSuccess(friendsParsed)
}

func Search(userSearchRequest models.UserSearchRequest, c chan utils.AppResponse) {
	valError := userSearchRequest.Validate()

	if len(valError) > 0 {
		c <- utils.ARValFail("", valError)
		return
	}

	var user models.User

	dbRes := utils.DB.Where("uuid = ?", userSearchRequest.UUID).First(&user)

	if dbRes.Error != nil {
		c <- utils.ARFailure("User not found.")
		return
	}

	c <- utils.ARSuccess(user.ToJSON())
}
