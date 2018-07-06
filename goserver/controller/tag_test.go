package controller

import "testing"

func TestFindTags(t *testing.T) {
  ret := findTags()
  if ret!= nil {
    t.Log("found it!")
    for idx,r := range ret{
      t.Logf("[%d] -------- %v \r\n",idx,r)
    }
  }else{
    t.Error("can't find !")
  }
}
