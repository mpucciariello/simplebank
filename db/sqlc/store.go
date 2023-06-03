package db

import (
	"context"
	"database/sql"
	"fmt"
)

type (
	Store struct {
		db *sql.DB
		*Queries
	}
	TransferTxParams struct {
		FromAccountID int64 `json:"from_account_id"`
		ToAccountID   int64 `json:"to_account_id"`
		Amount        int64 `json:"amount"`
	}
	TransferTxResult struct {
		Transfer      Transfer `json:"transfer"`
		FromAccountID Account  `json:"from_account_id"`
		ToAccountID   Account  `json:"to_account_id"`
		FromEntry     Entry    `json:"from_entry"`
		ToEntry       Entry    `json:"to_entry"`
	}
	BalanceTx struct {
		AccountID1 int64
		AccountID2 int64
		Amount1    int64
		Amount2    int64
	}
)

var txKey = struct{}{}

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

	txName := ctx.Value(txKey)

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		fmt.Println(txName, "create transfer")
		result.Transfer, err = s.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "create first entry")
		result.FromEntry, err = s.CreateEntry(ctx, CreateEntryParams{
			Amount:    -params.Amount,
			AccountID: params.FromAccountID,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "create second entry")
		result.ToEntry, err = s.CreateEntry(ctx, CreateEntryParams{
			Amount:    params.Amount,
			AccountID: params.ToAccountID,
		})

		if err != nil {
			return err
		}

		if params.FromAccountID < params.ToAccountID {
			result.FromAccountID, result.ToAccountID, err = modifyBalance(ctx, q, BalanceTx{
				AccountID1: params.FromAccountID,
				AccountID2: params.ToAccountID,
				Amount1:    -params.Amount,
				Amount2:    params.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccountID, result.FromAccountID, err = modifyBalance(ctx, q, BalanceTx{
				AccountID1: params.ToAccountID,
				AccountID2: params.FromAccountID,
				Amount1:    params.Amount,
				Amount2:    -params.Amount,
			})
			if err != nil {
				return err
			}
		}

		return nil

	})

	return result, err
}

func modifyBalance(ctx context.Context, q *Queries, balance BalanceTx) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: balance.Amount1,
		ID:     balance.AccountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: balance.Amount2,
		ID:     balance.AccountID2,
	})
	if err != nil {
		return
	}

	return account1, account2, nil
}
