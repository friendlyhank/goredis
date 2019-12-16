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

func TestRdsSetBit(t *testing.T){
	rs := GetRedisDefault()
	err := rs.SetBit("uid:100022:sign:2019-12-12",1,true)
	if err != nil{
		t.Errorf("%v",err)
		return
	}
}

func TestRdsGetBit(t *testing.T){
	num,err := GetRedisDefault().GetBit("uid:100022:sign:2019-12-12",1)
	if err != nil{
		t.Errorf("%v",err)
		return
	}
	t.Logf("%v",num)
}

func TestRdsBitCount(t *testing.T){
	count,err := GetRedisDefault().BitCount("uid:100022:sign:2019-12-12")
	if err != nil{
		t.Errorf("%v",err)
		return
	}
	t.Logf("%v",count)
}




