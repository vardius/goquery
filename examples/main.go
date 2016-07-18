package main

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vardius/goquery"
)

type (
	User struct {
		Id    int64  `column:"id"`
		Email string `column:"email"`
	}
)

func main() {
	var entities []*User

	reflectT := reflect.TypeOf(User{})
	builder := goquery.New(reflectT)
	conn, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/test")

	//SELECT
	data1, err := builder.Select().GetQuery().GetResults(conn)
	fmt.Println(data1, err)

	//SELECT SINGLE ROW
	data2, err := builder.Select().Where("id = ?").SetParameters(1).GetQuery().GetResult(conn)
	fmt.Println(data2, err)

	//REMOVE
	data3, err := builder.Delete().Where("id IN (?)").SetParameters([]int{1, 2}).GetQuery().Execute(conn)
	fmt.Println(data3, err)

	//ADD/UPDATE
	//builder will update this with id included
	entities = append(entities, &User{1, "test1@email.com"})
	entities = append(entities, &User{2, "test2@email.com"})
	//builder will add new row to table when id is not set
	entities = append(entities, &User{Email: "test2@email.com"})

	data4, err := builder.Save(entities).GetQuery().Execute(conn)
	fmt.Println(data4, err)
}
