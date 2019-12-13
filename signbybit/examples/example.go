package main

import (
	"fmt"
	_ "hank.com/goredis/foundation"
	"hank.com/goredis/foundation/rds"
)

func main(){
	err := rds.GetRedisDefault().SetBit("uid:100022:sign:2019-12-12",1,1)
	fmt.Println(err)

	num,_ := rds.GetRedisDefault().GetBit("uid:100022:sign:2019-12-12",1)
	fmt.Println(num)
}