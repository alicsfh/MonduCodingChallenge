package application

import (
	"context"
	"mondu-challenge-alihamedani/application/contracts"
	"mondu-challenge-alihamedani/domain"
	"mondu-challenge-alihamedani/domain/domain_services"
	util "mondu-challenge-alihamedani/infrastructure/utils"
	"time"
)

type (
	BankAccountService interface {
		AddAccount(ctx context.Context, model *contracts.CreateAccountModel) (contracts.AccountResult, error)
		DepositMoneyToAccount(ctx context.Context, model *contracts.CreateTransferModel) (domain.Account, error)
		WithdrawMoneyFromAccount(ctx context.Context, model *contracts.CreateTransferModel) (domain.Account, error)
		GetBalance(ctx context.Context, model *contracts.GetAccountBalanceModel) (contracts.AccountBalanceResult, error)
	}
	bankAccountService struct {
		Repository domain_services.IAccountRepository
	}
)

func (b bankAccountService) GetBalance(ctx context.Context, model *contracts.GetAccountBalanceModel) (contracts.AccountBalanceResult, error) {
	account, err := b.Repository.GetById(ctx, model.Id)

	if err != nil {
		return contracts.AccountBalanceResult{}, err
	}
	if account.ID() == 0 {
		return contracts.AccountBalanceResult{}, domain.AccountNotFoundError
	}

	result := contracts.AccountBalanceResult{
		FirstName: account.FirstName(),
		LastName:  account.LastName(),
		IBAN:      account.IBAN(),
		Balance:   account.Balance(),
	}
	return result, err
}

func (b bankAccountService) WithdrawMoneyFromAccount(ctx context.Context, model *contracts.CreateTransferModel) (domain.Account, error) {
	account, err := b.Repository.GetById(ctx, model.AccountOriginId)
	account.Withdraw(model.Amount)
	err = b.Repository.UpdateBalance(ctx, account)

	return account, err
}

func (b bankAccountService) DepositMoneyToAccount(ctx context.Context, model *contracts.CreateTransferModel) (domain.Account, error) {
	account, err := b.Repository.GetById(ctx, model.AccountOriginId)
	account.Deposit(model.Amount)
	err = b.Repository.UpdateBalance(ctx, account)

	return account, err
}

func (b bankAccountService) AddAccount(ctx context.Context, model *contracts.CreateAccountModel) (contracts.AccountResult, error) {
	id := util.NewId()
	account := domain.CreateAccount(id, model.FirstName, model.LastName, model.IBAN, model.Balance, time.Now())
	err := b.Repository.Create(ctx, account)

	result := contracts.AccountResult{
		Id:           account.ID(),
		FirstName:    account.FirstName(),
		LastName:     account.LastName(),
		IBAN:         account.IBAN(),
		Balance:      account.Balance(),
		CreationDate: account.CreationDate(),
	}
	return result, err

}

func NewAccountService(repository domain_services.IAccountRepository) BankAccountService {
	return &bankAccountService{Repository: repository}
}
