package controller

import "testing"

func TestLogin(t *testing.T) {
  r := login("boss","123456")
  if r.RetCode == "ok" {
    t.Log("success!")
    t.Logf("user name: %s \r\n",r.Name)
    t.Logf("user id: %s",r.Id)
  }else{
    t.Error("fail")
  }
}
