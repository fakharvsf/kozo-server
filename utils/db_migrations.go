package utils

import "rt-server/models"

func DBMigrate(noDB bool) {
	DB.AutoMigrate(&models.PersonalTask{}, &models.User{}, &models.FriendRequest{}, &models.Friend{}, &models.AssignedTask{})

	if noDB {
		DB.Migrator().CreateConstraint(&models.PersonalTask{}, "User")
		DB.Migrator().CreateConstraint(&models.PersonalTask{}, "fk_user_id")

		DB.Migrator().CreateConstraint(&models.FriendRequest{}, "Sender")
		DB.Migrator().CreateConstraint(&models.FriendRequest{}, "fk_sender_id")
		DB.Migrator().CreateConstraint(&models.FriendRequest{}, "Receiver")
		DB.Migrator().CreateConstraint(&models.FriendRequest{}, "fk_receiver_id")

		DB.Migrator().CreateConstraint(&models.Friend{}, "Sender")
		DB.Migrator().CreateConstraint(&models.Friend{}, "fk_sender_id")
		DB.Migrator().CreateConstraint(&models.Friend{}, "Receiver")
		DB.Migrator().CreateConstraint(&models.Friend{}, "fk_receiver_id")

		DB.Migrator().CreateConstraint(&models.AssignedTask{}, "Owner")
		DB.Migrator().CreateConstraint(&models.AssignedTask{}, "fk_owner_id")
		DB.Migrator().CreateConstraint(&models.AssignedTask{}, "AssignedTo")
		DB.Migrator().CreateConstraint(&models.AssignedTask{}, "fk_assigned_to_id")
		DB.Migrator().CreateConstraint(&models.AssignedTask{}, "PersonalTask")
		DB.Migrator().CreateConstraint(&models.AssignedTask{}, "fk_personal_task_id")
	}
}