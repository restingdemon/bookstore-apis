package helper

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/restingdemon/go-mysql-mux/pkg/models"
)


type SignedDetails struct{
	Email string
	First_name string
	Last_name string
	Uid string
	User_type string
	jwt.StandardClaims
}

var SECRET_KEY ="akshay"

func GenerateAllTokens(email string,firstName string,lastName string,userType string,uid string) (signedToken string , signedRefreshToken string, err error){
	claims := &SignedDetails{
		Email: email,
		First_name: firstName,
		Last_name: lastName,
		Uid: uid,
		User_type: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour*time.Duration(24)).Unix(),
		},
	}

	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token,err:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY))
	refeshToken ,err:=jwt.NewWithClaims(jwt.SigningMethodHS256,refreshclaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token,refeshToken,err
}


func UpdateAllTokens(signedToken string,signedRefreshToken string, userId string){

	var user models.User
	resp:=models.Db.Where("user_id=?",userId).First(&user)
	if resp.Error != nil{
		log.Panic(resp.Error)
		return
	}
	user.Token = signedToken
	user.Refresh_token = signedRefreshToken

	//filter:=models.Db.Where()

	models.Db.Model(user).Update(&user)

	// if err.Error!=nil {
	// 	log.Panic(err)
	// 	return
	// }
	return
}

func ValidateToken(signedToken string) (claims *SignedDetails,msg string){
	token,err:=jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY),nil
		},
	)

	if err != nil{
		msg = err.Error()
		return
	}

	claims,ok:=token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		return
	}
	return claims,msg
}

// ExtractTokenInfo extracts user type and user ID from the token
func ExtractTokenInfo(r *http.Request) (userType string, userID string, err error) {
	// Extract the token from the request headers
	token := r.Header.Get("token")
	if token == "" {
		return "", "", errors.New("No auth header provided")
	}

	// Parse the token to get claims
	claims, msg := ValidateToken(token)
	if msg != "" {
		return "", "", errors.New("Invalid token: " + msg)
	}

	return claims.User_type, claims.Uid, nil
}