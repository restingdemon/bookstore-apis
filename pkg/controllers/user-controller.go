package controllers

import (
	"encoding/json"

	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/restingdemon/go-mysql-mux/pkg/helper"
	"github.com/restingdemon/go-mysql-mux/pkg/models"
	"github.com/restingdemon/go-mysql-mux/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()
func HashPassword(password string) string{
	a,_:=bcrypt.GenerateFromPassword([]byte(password),14)

	return string(a)
}


func VerifyPassword( userpassword string, providedPass string)(bool,string){
	err:=bcrypt.CompareHashAndPassword([]byte(providedPass),[]byte(userpassword))
	check:= true
	msg:= ""
	if err!= nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check =false
	}
	return check,msg
}

func Signup(w http.ResponseWriter , r *http.Request){
	
	var user = &models.User{}
	utils.ParseBody(r,user)
	validationErr := validate.Struct(user)
	if validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validationErr.Error()))
		return
	}
	password:=HashPassword(*&user.Password)
	user.Password = password

	//check count for email and phone no to prevent multiple entries of same data
	err := user.BeforeCreate(models.Db)
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token,refreshToken,_ := helper.GenerateAllTokens(*&user.Email,*&user.First_name,*&user.Last_name,*&user.User_type,*&user.User_id)
	user.Token = token
	user.Refresh_token = refreshToken
	b := user.CreateUser()
	res , _ :=json.Marshal(b)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func Login(w http.ResponseWriter , r *http.Request){
	var user = &models.User{}
	utils.ParseBody(r,user)
	var Founduser models.User
	err:=models.Db.Where("email=?",user.Email).First(&Founduser)
	if err.Error!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(err.Error()))
		return
	}
	passwordIsValid,msg:=VerifyPassword(*&user.Password,*&Founduser.Password) 
	if !passwordIsValid {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	token,refreshtoken,_:=helper.GenerateAllTokens(*&Founduser.Email,*&Founduser.First_name,*&Founduser.Last_name,*&Founduser.User_type,*&Founduser.User_id)
	helper.UpdateAllTokens(token,refreshtoken,Founduser.User_id)
	err = models.Db.Where("user_id=?",Founduser.User_id).First(&Founduser)
	if err.Error!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg="yha pr error hua"
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(http.StatusOK)
	res,_:=json.Marshal(Founduser)
	w.Write(res)
}



func GetUsers(w http.ResponseWriter , r *http.Request){
	err:=helper.CheckUserType(r,"ADMIN")
	if err!=nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

}


// func GetUser(w http.ResponseWriter , r *http.Request){
// 	vars := mux.Vars(r)
// 	userId := vars["user_id"]

// 	err:= helper.MatchUserTypeToUid(r,userId)
// 	if err != nil {
// 		w.WriteHeader(http.StatusForbidden)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	userDetails,_ := models.GetUserById(userId)
// 	res,_:=json.Marshal(userDetails)
// 	w.Header().Set("Content-Type","pkglication/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(res)
// }
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]

	// Extract the user type and ID from the token
	tokenUserType, tokenUserID, err := helper.ExtractTokenInfo(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized: " + err.Error()))
		return
	}

	// Check if the requester is an admin
	if tokenUserType == "ADMIN" {
		// Admin can access any user's data
		userDetails, _ := models.GetUserById(userId)
		res, _ := json.Marshal(userDetails)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	// If the requester is a regular user, they can only access their own data
	if tokenUserID == userId {
		userDetails, _ := models.GetUserById(userId)
		res, _ := json.Marshal(userDetails)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	// If the requester is neither an admin nor the owner of the data, deny access
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Forbidden: User does not have access to this resource"))
}

