package usecase

import (
	"time"

	"github.com/figormartins/ms-it/domain"
	"github.com/figormartins/ms-it/dto"
)

type UseCaseTransaction struct {
	TransactionRepository domain.TransactionRepository
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) UseCaseTransaction {
	return UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (usecase UseCaseTransaction) ProcessTransaction(transactionDto dto.Transaction) (domain.Transaction, error) {
	creditCard := usecase.HydrateCreditCard(transactionDto)
	ccBalanceAndLimit, err := usecase.TransactionRepository.GetCreditCard(*creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}
	creditCard.Id = ccBalanceAndLimit.Id
	creditCard.Limit = ccBalanceAndLimit.Limit
	creditCard.Balance = ccBalanceAndLimit.Balance
	t := usecase.NewTransaction(transactionDto, ccBalanceAndLimit)
	t.ProcessAndValidate(creditCard)
	err = usecase.TransactionRepository.SaveTransaction(*t, *creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}

	return *t, nil
}

func (usecase UseCaseTransaction) HydrateCreditCard(transactionDto dto.Transaction) *domain.CreditCard {
	creditCard := domain.NewCreditCard()
	creditCard.Name = transactionDto.Name
	creditCard.Number = transactionDto.Number
	creditCard.ExpirationMonth = transactionDto.ExpirationMonth
	creditCard.ExpirationYear = transactionDto.ExpirationYear
	creditCard.CVV = transactionDto.CVV
	return creditCard
}

func (usecase UseCaseTransaction) NewTransaction(transaction dto.Transaction, cc domain.CreditCard) *domain.Transaction {
	t := domain.NewTransaction()
	t.CreditCardId = cc.Id
	t.Amount = transaction.Amount
	t.Store = transaction.Store
	t.Description = transaction.Description
	t.CreatedAt = time.Now()
	return t
}
