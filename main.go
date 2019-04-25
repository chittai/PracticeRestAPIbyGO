package main

import(
	"fmt"
	"time"
	"os"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/guregu/dynamo"
	"github.com/rs/xid"
	"github.com/joho/godotenv"
)

type User struct {
	UserID	xid.ID	`dynamo:"user_id"`
	Name	string	`dynamo:"name"`
	CreatedTime	time.Time	`dynamo:"created_time"`
}



func main(){
	godotenv.Load("./envfiles/develop.env")

	var apikey = os.Getenv("APIKEY")
	var secretkey = os.Getenv("SECRETKEY")

	var c = credentials.NewStaticCredentials(apikey, secretkey, "")

	var db = dynamo.New(session.New(), &aws.Config{
		Credentials: c,
		Region:	aws.String("us-east-1"),
	})
	
	var table = db.Table("Test")

	guid := xid.New()
	name := "taro"

	u := User{UserID: guid, Name: name, CreatedTime: time.Now().UTC()}
	fmt.Println(u)

	if err := table.Put(u).Run(); err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
}
