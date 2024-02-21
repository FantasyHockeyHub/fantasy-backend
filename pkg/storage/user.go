package storage

import (
	"database/sql"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
)

func (p *PostgresStorage) SignUp(u user.SignUpModel) error {
	u.ID = uuid.New()
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}
	err = p.CreateUserProfile(tx, u)
	if err != nil {
		return err
	}
	err = p.CreateUserData(tx, u)
	if err != nil {
		return err
	}
	err = p.CreateUserContacts(tx, u)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetUserInfo(userID uuid.UUID) (user.UserInfoModel, error) {
	var userInfo user.UserInfoModel

	err := p.db.Get(&userInfo, `SELECT user_profile.id, nickname, date_registration, photo_link, coins, uc.email FROM 
        user_profile INNER JOIN user_contacts uc on user_profile.id = uc.profile_id WHERE user_profile.id = $1;`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return userInfo, UserDoesNotExistError
		} else {
			return userInfo, err
		}
	}

	return userInfo, nil
}
