package simplepublish

import (
	"fmt"
	_ "github.com/friendlyhank/redis-use/foundation"
	"github.com/friendlyhank/goredis"
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
)

const(
	channelPrefix = "__redis"
)

//SimplePublish
type SimplePublish struct{
	//redis 源
	source *rds.RedisSource
	channel string
}

//channelName-
func channelName(name string)string{
	return fmt.Sprintf("%v_%v",channelPrefix,name)
}

//NewSimplePublish -
func NewSimplePublish()*SimplePublish{
	return &SimplePublish{
		source:rds.GetRedisDefault(),
		channel:channelName("simplepublish"),
	}
}

//Send-
func (s *SimplePublish)Send(message string){
	s.source.Publish(s.channel,message)
}

//StartRedisLoop-
func (s *SimplePublish)StartRedisLoop(){
	s.loopReceive()
	//go s.loopReceive()
}

func (s *SimplePublish)loopReceive(){
	pubSubConn := s.source.NewPubSubCoon()

	if err := pubSubConn.Subscribe(s.channel);err != nil{
		fmt.Println(err)
		return
	}

	for{
		switch n := pubSubConn.Receive().(type) {
		case redis.Message:
			log.Printf("pubSubConn Receiv Channel：%v;Pattern：%v;Data：%v",n.Channel,n.Pattern,string(n.Data))
		case redis.Subscription:
			log.Printf("pubSubConn Receiv Kind：%v;Channel：%v;Count：%v",n.Kind,n.Channel,n.Count)
		case redis.Pong:
			log.Printf("pubSubConn Receiv Data：%v",string(n.Data))
		}
	}

}

var (
	simplePublishOnce sync.Once
)

func Init(){
	//Do once
	simplePublishOnce.Do(func(){
		s := NewSimplePublish()
		s.StartRedisLoop()
	})
}


