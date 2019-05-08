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

	CreateSession()

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

func main(){
	router := mux.NewRouter()
	fmt.Println("Listening 8000 ...")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users", GetUsers).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
