package main

import (
	"fmt"

	"github.com/audi70r/go-struct2html/struct2html"
)

type User struct {
	Name  string
	Email string
}

func main() {
	users := []User{
		{Name: "John Doe", Email: "john.doe@example.com"},
		{Name: "Jane Smith", Email: "jane.smith@example.com"},
	}

	htmlTable, err := struct2html.StructToHTMLTable(users)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(htmlTable)
}
