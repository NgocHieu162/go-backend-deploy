package repository

import (
	"context"
	"errors"
	"go-backend/ent"
	"go-backend/internal/interfaces"
)

type keyTx struct{}

type UnitOfWorkRepository struct {
	entClient *ent.Client
}

func NewUnitOfWorkRepository(entClient *ent.Client) interfaces.UnitOfWorkRepository {
	return &UnitOfWorkRepository{
		entClient: entClient,
	}
}

// Do implements [interfaces.UnitOfWorkRepository].
func (u *UnitOfWorkRepository) Do(ctx context.Context, fn func(ctxTx context.Context) error) (err error) {
	tx, err := u.entClient.Tx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			// rollback
			errorRollback := tx.Rollback()
			if errorRollback != nil {
				err = errors.New("Rollback error")
			}
		} else {
			// commit
			errorCommit := tx.Commit()
			if errorCommit != nil {
				err = errors.New("commit error")
			}
		}
	}()

	ctxTx := context.WithValue(ctx, keyTx{}, tx.Client())

	err = fn(ctxTx)

	if err != nil {
		return err
	}

	return nil
}

func GetClientTx(ctx context.Context, client *ent.Client) *ent.Client {
	clientAny := ctx.Value(keyTx{})
	clientTx, ok := clientAny.(*ent.Client)
	if !ok {
		return client
	}

	return clientTx
}
