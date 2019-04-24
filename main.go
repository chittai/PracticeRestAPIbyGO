package main

import(
	"fmt"
	"time"
	"github.com/rs/xid"
)

type User struct {
	UserID	xid.ID	`dynamo:"user_id"`
	Name	string	`dynamo:"name"`
	CreatedTime	time.Time	`dynamo:"created_time"`
}

func main(){
	guid := xid.New()
	name := "taro"

	u := User{UserID: guid, Name: name, CreatedTime: time.Now().UTC()}
	fmt.Println(u)
}

