package repository

import (
	"database/sql"
	"errors"

	"github.com/figormartins/ms-it/domain"
)

type TransactionRepositoryDb struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepositoryDb {
	return &TransactionRepositoryDb{db: db}
}

func (transactionRepository *TransactionRepositoryDb) SaveTransaction(transaction domain.Transaction, creditCard domain.CreditCard) error {
	stmt, err := transactionRepository.db.Prepare(`insert into transactions(id, credit_card_id, amount, status, description, store, created_at)
																	values($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		transaction.Id,
		transaction.CreditCardId,
		transaction.Amount,
		transaction.Status,
		transaction.Description,
		transaction.Store,
		transaction.CreatedAt,
	)

	if err != nil {
		return err
	}

	if transaction.Status == "approved" {
		err = transactionRepository.UpdateBalance(creditCard)
		if err != nil {
			return err
		}
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepositoryDb) UpdateBalance(creditCard domain.CreditCard) error {
	_, err := t.db.Exec("update credit_cards set balance = $1 where id = $2",
		creditCard.Balance, creditCard.Id)
	if err != nil {
		return err
	}
	return nil
}

func (transactionRepository *TransactionRepositoryDb) CreateCreditCard(creditCard domain.CreditCard) error {
	stmt, err := transactionRepository.db.Prepare(`insert into credit_cards(id, name, number, expiration_month,expiration_year, CVV,balance, balance_limit) 
								values($1,$2,$3,$4,$5,$6,$7,$8)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		creditCard.Id,
		creditCard.Name,
		creditCard.Number,
		creditCard.ExpirationMonth,
		creditCard.ExpirationYear,
		creditCard.CVV,
		creditCard.Balance,
		creditCard.Limit,
	)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (transactionRepository *TransactionRepositoryDb) GetCreditCard(creditCard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard
	stmt, err := transactionRepository.db.Prepare("select id, balance, balance_limit from credit_cards where number=$1")
	if err != nil {
		return c, err
	}
	if err = stmt.QueryRow(creditCard.Number).Scan(&c.Id, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card does not exists")
	}
	return c, nil
}
