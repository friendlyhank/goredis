package rds

import "testing"

func SetInit(){
	Init()
}

func SetDestory(){
	rs := GetRedisDefault()
	rs.CloseConn(rs.dbpool.Get())
}

func TestMain(m *testing.M){
	SetInit()
	m.Run()
	SetDestory()
}

//TestRdsSetBit- 位的设置
func TestRdsSetBit(t *testing.T){
	rs := GetRedisDefault()
	err := rs.SetBit("uid:100022:sign:2019-12",1,true)
	if err != nil{
		t.Errorf("%v",err)
		return
	}
}

func TestRdsGetBit(t *testing.T){
	num,err := GetRedisDefault().GetBit("uid:100022:sign:2019-12",1)
	if err != nil{
		t.Errorf("%v",err)
		return
	}
	t.Logf("%v",num)
}

func TestRdsBitCount(t *testing.T){
	count,err := GetRedisDefault().BitCount("uid:100022:sign:2019-12")
	if err != nil{
		t.Errorf("%v",err)
		return
	}
	t.Logf("%v",count)
}

func TestRdsBitPos(t *testing.T){
	count,err := GetRedisDefault().BitPos("uid:100022:sign:2019-11",true)
	if err != nil{
		t.Errorf("%v",err)
		return
	}
	t.Logf("%v",count)
}

//TestIncr -数字增加
func TestIncr(t *testing.T){
	count := GetRedisDefault().Incr("10011_999")
	GetRedisDefault().SetExpire("10011_999",5)
	t.Logf("%v",count)
}

//TestGetUint64-获取数字
func TestGetUint64(t *testing.T){
	count,err := GetRedisDefault().GetUint64("10011_999")
	if err != nil{
		t.Errorf("%v",err)
		return
	}

	t.Logf("%v",count)

	ttl,err :=GetRedisDefault().TTL("10011_999")
	if err != nil{
		t.Errorf("%v",err)
		return
	}
	t.Logf("%v",ttl)
}




