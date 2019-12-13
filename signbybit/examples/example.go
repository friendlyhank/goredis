package main

import (
	"fmt"
	_ "hank.com/goredis/foundation"
	"hank.com/goredis/foundation/rds"
)

func main(){
	err := rds.GetRedisDefault().Ping()
	fmt.Println(err)
}