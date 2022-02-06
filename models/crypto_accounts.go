package models
//
//import (
//	mess "Crypto/json_message"
//	"github.com/dgrijalva/jwt-go"
//	"github.com/jinzhu/gorm"
//	"golang.org/x/crypto/bcrypt"
//	"os"
//	"strings"
//)
//
///*
//JWT claims struct
//*/
//type Token struct {
//	UserId uint
//	jwt.StandardClaims
//}
//
////a struct to rep user account
//type Crypto struct {
//	gorm.Model
//	Email    string `json:"email"`
//	Password string `json:"password"`
//	Token    string `json:"token";sql:"-"`
//}
//
////Validate incoming user details...
//func (account *Crypto) Validate() (map[string]interface{}, bool) {
//
//	if !strings.Contains(account.Email, "@") {
//		return mess.Message(false, "Email address is required"), false
//	}
//
//	if len(account.Password) < 6 {
//		return mess.Message(false, "Password is required"), false
//	}
//
//	//Email must be unique
//	temp := &Crypto{}
//
//	//check for errors and duplicate emails
//	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
//	if err != nil && err != gorm.ErrRecordNotFound {
//		return mess.Message(false, "Connection error. Please retry"), false
//	}
//	if temp.Email != "" {
//		return mess.Message(false, "Email address already in use by another user."), false
//	}
//
//	return mess.Message(false, "Requirement passed"), true
//}
//
//func (account *Crypto) Create() (map[string]interface{}) {
//
//	if resp, ok := account.Validate(); !ok {
//		return resp
//	}
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
//	account.Password = string(hashedPassword)
//
//	GetDB().Create(account)
//
//	if account.ID <= 0 {
//		return mess.Message(false, "Failed to create account, connection error.")
//	}
//
//	//Create new JWT token for the newly registered account
//	tk := &Token{UserId: account.ID}
//	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
//	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
//	account.Token = tokenString
//
//	account.Password = "" //delete password
//
//	response := mess.Message(true, "Account has been created")
//	response["account"] = account
//	return response
//}
//
//func Login(email, password string) (map[string]interface{}) {
//
//	account := &Crypto{}
//	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return mess.Message(false, "Email address not found")
//		}
//		return mess.Message(false, "Connection error. Please retry")
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
//	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
//		return mess.Message(false, "Invalid login credentials. Please try again")
//	}
//	//Worked! Logged In
//	account.Password = ""
//
//	//Create JWT token
//	tk := &Token{UserId: account.ID}
//	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
//	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
//	account.Token = tokenString //Store the token in the response
//
//	resp := mess.Message(true, "Logged In")
//	resp["account"] = account
//	return resp
//}
//
//func GetUser(u uint) *Crypto {
//
//	acc := &Crypto{}
//	GetDB().Table("accounts").Where("id = ?", u).First(acc)
//	if acc.Email == "" { //User not found!
//		return nil
//	}
//
//	acc.Password = ""
//	return acc
//}