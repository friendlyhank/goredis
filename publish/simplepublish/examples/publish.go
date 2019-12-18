package main

import "github.com/friendlyhank/goredis/publish/simplepublish"

func main(){
	s := simplepublish.NewSimplePublish()
	s.Send("hello world")
}
