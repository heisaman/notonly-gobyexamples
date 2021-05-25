package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "user:password@tcp(localhost:3360)/testDB")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	ctx := context.Background()
	osType := "windows"
	environment := "生产"
	networkType := "classic"
	results, err := db.QueryContext(ctx, "SELECT public_ip_address FROM create_ecs_instance where os_type=? and environment=? and instance_network_type=?", osType, environment, networkType)
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var ip string
		err = results.Scan(&ip)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(ip)
	}
}
