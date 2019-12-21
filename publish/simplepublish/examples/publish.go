package main

import "github.com/friendlyhank/redis-use/publish/simplepublish"

func main(){
	s := simplepublish.NewSimplePublish()
	s.Send("hello world")
}
