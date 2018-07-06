package util

import (
  "testing"
  "github.com/magicgravity/vgblog/goserver/model"

  "time"
)

func TestCopyProperty(t *testing.T) {
  u := model.User{}
  u.Name="testadmin"
  u.Salt="122222"
  u.Password = "33333"

  u2 := model.User{}
  if !CopyProperty(&u,&u2) {
    t.Error("copy property fail!")
  }else{
    t.Log("copy property successful !")
    t.Logf("new user >> %v ",u2)
  }
}


func TestCopyProperty2(t *testing.T) {
  type User1 struct {
    UserId string
    Name string
    Grade int
    Sex uint
    RankSeq uint64
    Birth time.Time
    Amount float64
  }

  type User2 struct {
    Amount string
    UserId string
    Name string
    Grade string
    Birth string
    Sex string
    RankSeq string
  }

  u1 := User1{}
  u1.Name="哈哈"
  u1.Birth= time.Now()
  u1.Grade=3
  u1.Sex=1
  u1.RankSeq=302133131
  u1.UserId= "234k23123ppp"
  u1.Amount = 923923.88

  u2 := User2{}
  if CopyProperty(&u1,&u2) {
    t.Log("copy property successful !")
    t.Logf("user2 >> %v ",u2)
  }else{
    t.Error("copy property fail!")
  }
}
