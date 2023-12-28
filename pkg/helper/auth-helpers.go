package helper

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func CheckUserType(r *http.Request,  role string)(err error){
	vars := mux.Vars(r)
	userType := vars["user_type"] 
	err = nil
	if userType != role {
		err = errors.New("Unauthorized to access this resource (2)")
		return err
	}
	return err
}

func MatchUserTypeToUid(r *http.Request,userId string) (err error){
	vars := mux.Vars(r)
	userType := vars["user_type"]
	uId := vars["user_id"]
	err =nil	

	if userType == "USER" && uId != userId {
		err = errors.New("Unauthorised to access this resouce (1)")
		return err
	}

	err = CheckUserType(r,userType)
	return err
}
