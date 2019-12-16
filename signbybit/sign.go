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
func (us *UserSignDate)DoSign(uid int64,localTime *time.Time,status bool){
	offset := localTime.Day()
	rds.GetRedisDefault().SetBit(GetUserSignKey(uid,localTime),int64(offset),status)
}

//CheckSign-检查用户是否签到
func (us *UserSignDate)CheckSign(uid int64,localTime *time.Time)bool{
	offset := localTime.Day()
	value,_ := rds.GetRedisDefault().GetBit(GetUserSignKey(uid,localTime),int64(offset))
	return value > 0
}

//GetSignCount -获取用户签到的次数
func (us *UserSignDate)GetSignCount(uid int64,localTime *time.Time)(int,error){
	return rds.GetRedisDefault().BitCount(GetUserSignKey(uid,localTime))
}

func GetUserSignKey(uid int64,localTime *time.Time)string{
	return fmt.Sprintf("uid:%v:sign:%v",uid,localTime.Format("2006-01-02"))
}
