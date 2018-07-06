package main

import (
  "github.com/magicgravity/vgblog/goserver/model"
  "github.com/magicgravity/vgblog/goserver/util"
  "fmt"
)

func main(){
  u := model.User{}
  u.Name="testadmin"
  u.Salt="122222"
  u.Password = "33333"

  u2 := model.User{}

  if util.CopyProperty(&u,&u2) {
    fmt.Println("copy ok!")
    fmt.Printf("new user >> %v \r\n",u2)
  }else{
    fmt.Println("copy fail")
  }

  u3 := SomeUser{}
  if util.CopyProperty(&u,&u3) {
    fmt.Println("copy ok!")
    fmt.Printf("new some user >> %v \r\n",u3)
  }else{
    fmt.Println("copy fail")
  }
}


type SomeUser struct {
  UserId string
  Name string
  Password string
  Salt string
  Sex bool
  IdNo string
}
