package storage

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
)

func (p *PostgresStorage) SignUp(ctx context.Context, u user.SignUpModel) error {
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
