package signbybit

import (
	"fmt"
	"time"
	_ "hank.com/goredis/foundation"
	"hank.com/goredis/foundation/rds"
)

//UserSignDate -用户签到
type UserSignDate struct{

}

//DoSign -设置签到
func (us *UserSignDate)DoSign(uid int64,localTime *time.Time){
	key := GetUserSignKey(uid,localTime)
	rds.GetRedisDefault().SetBit(key,1,1)
}

func GetUserSignKey(uid int64,localTime *time.Time)string{
	return fmt.Sprintf("uid:%v:sign:%v",uid,localTime.Format("2006-01-02"))
}
