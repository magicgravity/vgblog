package controller

import (
  "gitee.com/johng/gf/g/net/ghttp"
  "github.com/magicgravity/vgblog/goserver/db"
  "gopkg.in/mgo.v2/bson"
  "github.com/magicgravity/vgblog/goserver/model"
  "github.com/magicgravity/vgblog/goserver/pojo"
  "encoding/json"
  "github.com/magicgravity/vgblog/goserver/util"
)

func init(){
  ghttp.GetServer().BindHandler("get:/api/tags",     FindTags)
}


func FindTags(r *ghttp.Request){
  tags := findTags()
  if tags==nil {
    util.FailRet(500,"",r)
  }else{
    respData,err := json.Marshal(tags)
    if err!= nil {
      util.FailRet(500,"",r)
    }else{
      util.SucceedRet(string(respData),r)
    }
  }
}

//func FindTags(r *ghttp.Request){
//  articles := findTags()
//  if articles==nil {
//    r.Response.WriteJson("{}")
//  }else{
//    resps := make([]pojo.TagResp,len(articles))
//    for idx,a := range articles{
//      rr := pojo.TagResp{}
//      rr.Id = a.Id_.Hex()
//      rr.Content = a.Content
//      rr.Title = a.Title
//      rr.Aid = a.Aid
//      rr.CommentN = a.CommentN
//      rr.IsPublish = a.IsPublish
//      copy(rr.Tags,a.Tags)
//      rr.Date = util.FormatDate(a.Date,"YYYYMMDD HH:MI:SS")
//
//      resps[idx] = rr
//    }
//    respData,err := json.Marshal(resps)
//    if err!= nil {
//      r.Response.WriteHeader(500)
//    }else{
//      r.Response.WriteHeader(200)
//      r.Response.WriteJson(string(respData))
//    }
//  }
//}

//func findTags()[]model.Article{
//  sess := db.GetMongoConnFromPool()
//  defer db.ReturnMongoConn(sess)
//  cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
//
//  rets := make([]model.Article,0)
//
//  err := cu.Find(bson.M{"isPublish":true}).Distinct("tags",&rets)
//  if err!= nil {
//    return nil
//  }else{
//    return rets
//  }
//}



func findTags()[]string{
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)

  rets := make([]string,0)

  err := cu.Find(bson.M{"isPublish":true}).Distinct("tags",&rets)
  if err!= nil {
    return nil
  }else{
    return rets
  }
}
