package transaction

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}
type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {

	//cek authorization
	campaign, err := s.campaignRepository.FindById(input.ID)

	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserId != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the campaign")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)

	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {

	uuid := uuid.New()

	//buat objek transaction
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	transaction.Code = fmt.Sprintf("ORDER-%s", uuid)

	newTransaction, err := s.repository.Save(transaction)

	if err != nil {
		return newTransaction, err
	}
	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Code:   newTransaction.Code,
		Amount: newTransaction.Amount,
	}

	paymentUrl, err := s.paymentService.GetPaymentUrl(paymentTransaction, input.User)

	if err != nil {
		log.Print("Error di sini")
		return newTransaction, err
	}

	newTransaction.PaymentUrl = paymentUrl
	// fmt.Println("Cek")
	// fmt.Println(paymentUrl)
	newTransaction, err = s.repository.Update(newTransaction)

	if err != nil {
		log.Print("Error di sana")
		return newTransaction, err
	}
	return newTransaction, nil
}
