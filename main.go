package main

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/valyala/fasthttp"

	//"github.com/go-chi/chi"
	//"github.com/go-chi/render"
	"os"
	"strings"
	"time"

	//"github.com/rs/cors"
	"log"
	//_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Crypto struct {
	gorm.Model
	Name  string
	Price float64
}

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Account struct {
	gorm.Model
	Email    string
	Password []byte
	Token    string
}

var db *gorm.DB
var err error
var was_logged_in = false

var (

	cryptoValues = []Crypto{
		{Name: "TABOO TOKEN", Price: 0.855},
		{Name: "Rainbow Token",  Price: 0.0000004251},
		{Name: "Ariva",  Price: 0.0678},
	}

)

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome to Crypto Service!\n")
}

func Show(ctx *fasthttp.RequestCtx) {
	if was_logged_in {
		var crypto []Crypto
		db.Find(&crypto)
		json.NewEncoder(ctx).Encode(&crypto)
	} else {
		fmt.Fprint(ctx, "Please, log in!\n")
	}
}

func DeleteValue(ctx *fasthttp.RequestCtx) {
	if was_logged_in {
		db, err = gorm.Open( "postgres", "host=127.0.0.1 port=5432 user=postgres dbname=crypto_values sslmode=disable password=gfhjkm")
		name := ctx.UserValue("name").(string)
		//id_int, err := strconv.Atoi(id)
		//fmt.Println(id_int, err, reflect.TypeOn(id_int))
		db.Where("Name = ?", name).Delete(&Crypto{})
		fmt.Fprint(ctx, "Was deleted successfully!\n")
	} else {
		fmt.Fprint(ctx, "Please, log in!\n")
	}
}


func AddValue(ctx *fasthttp.RequestCtx) {
	if was_logged_in {
		//db := dbConn()
		db, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=crypto_values sslmode=disable password=gfhjkm")
		//fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
		name := ctx.UserValue("name").(string)
		//name := string(ctx.FormValue("name"))
		value := Crypto{Name: name, Price: 0.9999}
		db.Select("Name", "Price").Create(&value)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprint(ctx, "Was added successfully!\n")
	} else {
		fmt.Fprint(ctx, "Please, log in!\n")
	}
}

//func New(ctx *fasthttp.RequestCtx) {
//	tmpl.ExecuteTemplate(ctx, "New", nil)
//}

//func Update(w http.ResponseWriter, r *http.Request) {
//	db := dbConn()
//	if r.Method == "POST" {
//		name := r.FormValue("name")
//		city := r.FormValue("city")
//		id := r.FormValue("uid")
//		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
//		if err != nil {
//			panic(err.Error())
//		}
//		insForm.Exec(name, city, id)
//		log.Println("UPDATE: Name: " + name + " | City: " + city)
//	}
//	defer db.Close()
//	http.Redirect(w, r, "/", 301)
//}
//

//check email, password and duplicate
func (account *Account) Validate() (bool) {

	if !strings.Contains(account.Email, "@") {
		fmt.Println("Bad email!")
		return false
	}

	if len(account.Password) < 6 {
		fmt.Println("Bad password!")
		return false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := db.Table("accounts").Where("email = ?", account.Email).Find(&temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("Bad ErrRecordNotFound!")
		return false
	}
	if temp.Email != "" {
		fmt.Println("Email != \"\"!")
		return false
	}

	return true
}

func CreateAccount(ctx *fasthttp.RequestCtx) {
	db, err = gorm.Open( "postgres", "host=127.0.0.1 port=5432 user=postgres dbname=accounts sslmode=disable password=gfhjkm")
	if err != nil{
		fmt.Fprint(ctx, "Can't open account database!\n")
	}
	name := ctx.UserValue("name").(string)
	password := ctx.UserValue("password").(string)
	account:=Account{Email: name, Password: []byte(password), Token: ""}
	if !account.Validate(){
		fmt.Fprint(ctx, "Bad credentials!\n")
		return
	}

	//Create account
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = hashedPassword
	//Create new JWT token for the newly registered account
	tk := Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("gfhjkm"))
	account.Token = tokenString
	db.Select("Email", "Password", "Token").Create(&account)
	if account.ID <= 0 {
		fmt.Fprint(ctx, "Failed to create account, connection error!\n")
		fmt.Println(string(account.ID))
	}
	//print account database
	var acc []Account
	db.Find(&acc).Create(account)
	json.NewEncoder(ctx).Encode(&acc)

}

func DeleteAccount(ctx *fasthttp.RequestCtx){
	db, err = gorm.Open( "postgres", "host=127.0.0.1 port=5432 user=postgres dbname=accounts sslmode=disable password=gfhjkm")
	email := ctx.UserValue("email").(string)
	//id_int, err := strconv.Atoi(id)
	//fmt.Println(id_int, err, reflect.TypeOn(id_int))
	db.Where("Email = ?", email).Delete(&Account{})
	fmt.Fprint(ctx, "Account was deleted successfully!\n")
}

func Login(ctx *fasthttp.RequestCtx) {
//func Login(email, password string, ctx *fasthttp.RequestCtx) (map[string]interface{}) {
	got_email := ctx.UserValue("email").(string)
	got_password := ctx.UserValue("password").(string)
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(got_password), bcrypt.DefaultCost)
	account := &Account{}
	db, err = gorm.Open( "postgres", "host=127.0.0.1 port=5432 user=postgres dbname=accounts sslmode=disable password=gfhjkm")
	//err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	db.Where("Email = ?", got_email).Find(&account)
	//db.Where("Email = @name", sql.Named("Email", "@new_NatatataIvan")).Find(&account)
	if account.Email ==""{
		fmt.Fprint(ctx,"Bad login")
		return
	}
	err = bcrypt.CompareHashAndPassword(account.Password, []byte(got_password))
	if err != nil || err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		fmt.Fprint(ctx,"Invalid login credentials. Please try again")
		return
	}
	//account.Password = ""

	//Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	fmt.Fprint(ctx, "Logged In")
	was_logged_in = true
	//resp["account"] = account
	//return resp
}

func main() {
	db, err = gorm.Open( "postgres", "host=127.0.0.1 port=5432 user=postgres dbname=crypto_values sslmode=disable password=gfhjkm")

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Crypto{})
	db, err = gorm.Open( "postgres", "host=127.0.0.1 port=5432 user=postgres dbname=accounts sslmode=disable password=gfhjkm")
	db.AutoMigrate(&Account{})
	//for index := range cryptoValues {
	//	db.Create(&cryptoValues[index])
	//}
	router := fasthttprouter.New()
	log.Println("Server started on: http://localhost:8080")
	router.GET("/", Index)
	//http.HandleFunc("/edit", Edit)
	//http.HandleFunc("/update", Update)
	router.GET("/add_acc/:name/:password", CreateAccount)
	router.GET("/delete_acc/:email", DeleteAccount)
	router.GET("/login/:email/:password", Login)
	//router.POST("/login", Login)
	router.GET("/show", Show)
	router.GET("/add/:name", AddValue)
	router.GET("/delete/:name", DeleteValue)

	//handler := cors.Default().Handler(router)
	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}


type User struct {
	ID uint64            `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone string `json:"phone"`
}
var user = User{
	ID:            1,
	Username: "Nat",
	Password: "Password",
	Phone: "49123454322", //this is a random number
}


func CreateToken(userId uint64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}