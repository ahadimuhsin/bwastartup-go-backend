package payment

import (
	"bwastartup/user"
	"log"
	"os"

	// "strconv"

	"github.com/joho/godotenv"
	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
	
}

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
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

