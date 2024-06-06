# RediSQL 
_This is a personal pet project, not a maintained library or a production ready software._


# Examples
```go
package main

import (
	"fmt"
	"github.com/juricake/redisql/pkg/redisql"
)

func main() {
	// client init
	client, _ := redisql.NewClient(redisql.Options{
		Host: "localhost",
		Port: 6379,
	})

	// table init
	if err := client.CreateTable("users", user{}); err != nil {
		panic(err)
	}
	
	// insert the data
	sample := user{ID: 1, Name: "jon.doe", Age: 42}
	if err := client.InsertInto("users", sample); err != nil {
		panic(err)
	}

	// extract the data
	var results []user
	if err := client.SelectFrom("users", &results); err != nil {
		panic(err)
	}

	for _, res := range results {
		fmt.Println(res.Name)
	}
}

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}


```
