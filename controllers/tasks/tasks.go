package tasks

import (
	"kozo/models"
	"kozo/utils"
)

func Create(ID uint, personalTaskCreate models.PersonalTaskCreate, c chan utils.AppResponse) {
	valError := personalTaskCreate.Validate()

	if len(valError) > 0 {
		c <- utils.ARValFail("", valError)
		return
	}

	var user models.User

	_ = utils.DB.First(&user, ID)

	personalTask := models.PersonalTask{
		UserID: ID,
		Title: personalTaskCreate.Title,
		Description: personalTaskCreate.Description,
		Status: models.PersonalTaskStatusesEnums["created"],
	}

	_ = utils.DB.Create(&personalTask)
	personalTask.User = user

	res := utils.ARSuccess(personalTask.ToJSON())

	c <- res
}

func ReadOne(ID uint, personalTaskID uint, c chan utils.AppResponse) {
	personalTask := models.PersonalTask{}

	dbRes := utils.DB.Preload("User").First(&personalTask, personalTaskID)

	if dbRes.Error != nil {
		c <- utils.ARFailure("Personal task not found.")
		return
	}

	if personalTask.UserID != ID {
		c <- utils.ARFailure("You are not authorized to read this personal task.")
		return
	}

	res := utils.ARSuccess(personalTask.ToJSON())

	c <- res
}

func ReadAll(ID uint, c chan utils.AppResponse) {
	personalTasks := []models.PersonalTask{}

	utils.DB.Preload("User").Where("user_id = ?", ID).Find(&personalTasks)

	personalTasksParsed := []models.PersonalTaskJSON{}

	for i := 0; i < len(personalTasks); i++ {
		obj := personalTasks[i]
		personalTasksParsed = append(personalTasksParsed, obj.ToJSON())
	}

	res := utils.ARSuccess(personalTasksParsed)

	c <- res
}

func ReadPersonalAndAssigned(ID uint, c chan utils.AppResponse) {
	personalTasks := []models.PersonalTask{}
	utils.DB.Preload("User").Where("user_id = ?", ID).Find(&personalTasks)
	personalTasksParsed := []models.PersonalTaskJSON{}
	for i := 0; i < len(personalTasks); i++ {
		obj := personalTasks[i]
		personalTasksParsed = append(personalTasksParsed, obj.ToJSON())
	}

	assignedTasks := []models.AssignedTask{}
	utils.DB.Preload("Owner").Preload("AssignedTo").Preload("PersonalTask").Where("assigned_to_id = ?", ID).Find(&assignedTasks)
	assignedTasksParsed := []models.AssignedTaskJSON{}
	for i := 0; i < len(assignedTasks); i++ {
		obj := assignedTasks[i]
		obj.PersonalTask.User = obj.Owner
		assignedTasksParsed = append(assignedTasksParsed, obj.ToJSON())
	}

	personalAndAssignedTasks := models.PersonalAndAssignedTasks{}

	personalAndAssignedTasks.PersonalTasks = personalTasksParsed
	personalAndAssignedTasks.AssignedTasks = assignedTasksParsed

	res := utils.ARSuccess(personalAndAssignedTasks)

	c <- res
}

func Update(ID uint, personalTaskID uint64, personalTaskUpdate models.PersonalTaskUpdate, c chan utils.AppResponse) {
	valError := personalTaskUpdate.Validate()

	if len(valError) > 0 {
		c <- utils.ARValFail("", valError)
		return
	}

	personalTask := models.PersonalTask{}

	dbRes := utils.DB.Preload("User").First(&personalTask, uint(personalTaskID))

	if dbRes.Error != nil {
		c <- utils.ARFailure("Personal task not found.")
		return
	}

	if personalTask.UserID != ID {
		c <- utils.ARFailure("You are not authorized to edit this task.")
		return
	}

	if personalTaskUpdate.Status != "" {
		personalTask.Status = personalTaskUpdate.Status
	}

	if personalTaskUpdate.Title != "" {
		personalTask.Title = personalTaskUpdate.Title
	}

	if personalTaskUpdate.Description != "" {
		personalTask.Description = personalTaskUpdate.Description
	}

	utils.DB.Save(&personalTask)

	res := utils.ARSuccess(personalTask.ToJSON())

	c <- res
}

func Delete(ID uint64, c chan utils.AppResponse) {
	var personalTask models.PersonalTask

	dbRes := utils.DB.Preload("User").First(&personalTask, uint(ID))

	if dbRes.Error != nil {
		c <- utils.ARFailure("Personal task not found.")
	}

	utils.DB.Delete(&personalTask)

	res := utils.ARSuccess(personalTask.ToJSON())

	c <- res
}

func Assign(ID uint, assignPersonalTask models.AssignPersonalTask, c chan utils.AppResponse) {
	valError := assignPersonalTask.Validate()

	if len(valError) > 0 {
		c <- utils.ARValFail("", valError)
		return
	}

	var owner models.User

	_ = utils.DB.First(&owner, ID)

	var assignedTo models.User

	dbResAssignee := utils.DB.First(&assignedTo, assignPersonalTask.AssignedToID)

	if dbResAssignee.Error != nil {
		c <- utils.ARFailure("Assigned to user not found",)
		return
	}

	var personalTask models.PersonalTask

	dbResPersonalTask := utils.DB.Preload("User").First(&personalTask, assignPersonalTask.PersonalTaskID)

	if dbResPersonalTask.Error != nil {
		c <- utils.ARFailure("Personal task not found",)
		return
	}

	assignedTaskRes := utils.DB.Where("assigned_to_id = ?", assignedTo.ID).First(&models.AssignedTask{})
	assignedTaskExists := assignedTaskRes.Error == nil

	if assignedTaskExists {
		c <- utils.ARFailure("This task has already been assigned to this user.")
		return
	}

	friendRes := utils.DB.Where("(sender_id = ? AND receiver_id = ?) OR (receiver_id = ? AND sender_id = ?)", owner.ID, assignedTo.ID, owner.ID, assignedTo.ID).First(&models.Friend{})
	areFriends := friendRes.Error == nil

	if !areFriends {
		c <- utils.ARFailure("You are not friends with the assigned to user.")
		return
	}

	if personalTask.UserID != ID {
		c <- utils.ARFailure("You are not owner of this task.")
		return
	}

	if personalTask.UserID == assignPersonalTask.AssignedToID {
		c <- utils.ARFailure("Users can't assign tasks to themselves.")
		return
	}

	assignedTask := models.AssignedTask{
		OwnerID: ID,
		AsssignedToID: assignedTo.ID,
		PersonalTaskID: personalTask.ID,
	}

	utils.DB.Create(&assignedTask)

	assignedTask.PersonalTask.Status = models.PersonalTaskStatusesEnums["assigned"]
	utils.DB.Save(&assignedTask)

	assignedTask.Owner = owner
	assignedTask.AssignedTo = assignedTo
	assignedTask.PersonalTask = personalTask

	res := utils.ARSuccess(assignedTask.ToJSON())

	c <- res
}

func ReadAssignedTasks(ID uint, c chan utils.AppResponse) {
	assignedTasks := []models.AssignedTask{}

	utils.DB.Preload("Owner").Preload("AssignedTo").Preload("PersonalTask").Where("assigned_to_id = ?", ID).Find(&assignedTasks)

	personalTasksParsed := []models.PersonalTaskJSON{}

	for i := 0; i < len(assignedTasks); i++ {
		obj := assignedTasks[i]
		obj.PersonalTask.User = obj.Owner
		personalTasksParsed = append(personalTasksParsed, obj.PersonalTask.ToJSON())
	}

	res := utils.ARSuccess(personalTasksParsed)

	c <- res
}
