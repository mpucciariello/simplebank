package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
	*Queries
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer      Transfer `json:"transfer"`
	FromAccountID Account  `json:"from_account_id"`
	ToAccountID   Account  `json:"to_account_id"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx receives a function as a parameter and executes it within the database transaction
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil) //the second parameter of the function defines the level of isolation. nil equals to the default value
	if err != nil {
		return err
	}

	q := New(tx) // New accepts any value within the DBTX interface
	err = fn(q)
	if err != nil {
		// rollback
		if rbErr := tx.Rollback(); rbErr != nil {
			// append error
			return fmt.Errorf("transaction error: %s, rollback error: %s", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTx executes a query performing all the necessary db transactions involved in a transfer
// It creates the transfer register, creates the account entries and updates the balance in both accounts within a single database transaction
func (s *Store) TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = s.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = s.CreateEntry(ctx, CreateEntryParams{
			Amount:    -params.Amount,
			AccountID: params.FromAccountID,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = s.CreateEntry(ctx, CreateEntryParams{
			Amount:    params.Amount,
			AccountID: params.ToAccountID,
		})

		if err != nil {
			return err
		}

		return nil

	})

	return result, err
	
}
