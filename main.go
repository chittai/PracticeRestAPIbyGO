package main

import(
	"fmt"
	"time"
	"os"
	"net/http"
	"log"
	"encoding/json"

	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/guregu/dynamo"
	"github.com/rs/xid"
	"github.com/joho/godotenv"
)

var apikey string
var secretkey string
var c *credentials.Credentials
var db *dynamo.DB
var table dynamo.Table

type User struct {
	UserID	xid.ID	`dynamo:"user_id"`
	Name	string	`dynamo:"name"`
	CreatedTime	time.Time	`dynamo:"created_time"`
}

func CreateSession() {
	godotenv.Load("./envfiles/develop.env")
	apikey = os.Getenv("APIKEY")
	secretkey = os.Getenv("SECRETKEY")
	
	c = credentials.NewStaticCredentials(apikey, secretkey, "")

	db = dynamo.New(session.New(), &aws.Config{
		Credentials: c,
		Region:	aws.String("us-east-1"),
	})
	table = db.Table("Test")
}

func CreateUser(w http.ResponseWriter, r *http.Request){
	godotenv.Load("./envfiles/develop.env")

	guid := xid.New()
	name := r.URL.Query().Get("name")
	
	u := User{UserID: guid, Name: name, CreatedTime: time.Now().UTC()}
	fmt.Println(u)
	if err := table.Put(u).Run(); err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	var users []User
	err := table.Scan().All(&users)
	if err != nil {
		fmt.Println("err")
		panic(err.Error())
	}

	for i := range users {
		fmt.Println(users[i])
	}
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var user []User
	err := table.Get("user_id", userID).All(&user)
	if err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
	fmt.Println(user)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	q := r.URL.Query().Get("name")

	userID := params["id"]

	err := table.Update("user_id", userID).Set("name", q).Run()
	if err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]
	table.Delete("user_id", userID).Run()
}


func main(){

	CreateSession()
	router := mux.NewRouter()
	fmt.Println("Listening 8000 ...")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", UpdateUser).Methods("POST")
	router.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
