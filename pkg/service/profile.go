package service

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
	"log"
)

func (s *UserService) GetUserInfo(userID uuid.UUID) (user.UserInfoModel, error) {
	userInfo, err := s.storage.GetUserInfo(userID)
	if err != nil {
		log.Println("Service. GetUserInfo:", err)
		return userInfo, err
	}
	return userInfo, nil
}

func (s *UserService) CheckUserDataExists(inp user.UserExistsDataInput) error {
	var exists bool
	var err error

	switch {
	case inp.Email != "":
		exists, err = s.storage.CheckEmailExists(inp.Email)
		if err != nil {
			log.Println("Service. CheckEmailExists:", err)
			return err
		}
	case inp.Nickname != "":
		exists, err = s.storage.CheckNicknameExists(inp.Nickname)
		if err != nil {
			log.Println("Service. CheckNicknameExists:", err)
			return err
		}
	}

	if exists {
		return nil
	}

	return UserDoesNotExistError
}

func (s *UserService) DeleteProfile(userID uuid.UUID) error {

	err := s.storage.DeleteProfile(userID)
	if err != nil {
		log.Println("Service. DeleteProfile:", err)
		return err
	}

	return nil
}
