package auth

import (
	// "fmt"
	"errors"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Service interface {
	GenerateToken(userId int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewService() *jwtService {
	return &jwtService{}
}
var SECRET_KEY="THE MOMENT PEOPLE COME TO KNOW LOVE, THEY RUN THE RISK OF CARRYING HATE"
// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// fmt.Println(key)
	return os.Getenv(key)
}

func (s *jwtService) GenerateToken(userId int) (string, error) {
	//siapkan data yg disisipkan ke dalam token
	claim := jwt.MapClaims{}
	claim["user_id"] = userId

	//generate token
	//parameter satu algoritmanya, parameter kedua adalah claim/payloadnya
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(goDotEnvVariable("SECRET_KEY")))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error){
	//parse token
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error){
		 _, ok := token.Method.(*jwt.SigningMethodHMAC)
		 //jika methodnya bukan HMAC
		 if !ok{
			return nil, errors.New("Invalid token")
		 }
		 return []byte(goDotEnvVariable("SECRET_KEY")), nil
	})

	if err != nil{
		return token, err
	}

	return token, nil
}
