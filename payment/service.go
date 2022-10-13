package payment

import (
	"bwastartup/campaign"
	"bwastartup/transaction"
	"bwastartup/user"
	"log"
	"os"

	// "strconv"

	"github.com/joho/godotenv"
	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
	transactionRepository transaction.Repository
	campaignRepository    campaign.Repository
}

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
	ProcessPayment(input transaction.TransactionNotificationInput) error
}

func NewService(transactionRepository transaction.Repository, campaignRepository campaign.Repository) *service {
	return &service{}
}

func (s *service) GetPaymentUrl(transaction Transaction, user user.User) (string, error) {
	//env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ServerKey := os.Getenv("SERVER_KEY")
	ClientKey := os.Getenv("CLIENT_KEY")

	//midtrans init
	midclient := midtrans.NewClient()
	midclient.ServerKey = ServerKey
	midclient.ClientKey = ClientKey
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			// OrderID: "ORDER-ID-"+strconv.Itoa(transaction.ID),
			OrderID:  transaction.Code,
			GrossAmt: int64(transaction.Amount),
		},
	}

	// log.Println("GetToken:")
	snapTokenResp, err := snapGateway.GetToken(snapReq)

	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}

func (s *service) ProcessPayment(input transaction.TransactionNotificationInput) error {
	//ambil data transaksi berdasarkan order id
	transaction_id := input.OrderID

	transaction, err := s.transactionRepository.GetByOrderID(transaction_id)

	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	//update data status transaction
	updatedTransaction, err := s.transactionRepository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindById(updatedTransaction.CampaignID)

	if err != nil {
		return err
	}

	//update backer count dan current amount
	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)

		if err != nil {
			return err
		}
	}

	return nil
}
