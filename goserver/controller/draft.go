package controller

import (
  "gitee.com/johng/gf/g/net/ghttp"
  "github.com/magicgravity/vgblog/goserver/model"
  "time"
  "github.com/magicgravity/vgblog/goserver/util"
  "strconv"
  "github.com/magicgravity/vgblog/goserver/db"
  "gopkg.in/mgo.v2/bson"
  "github.com/magicgravity/vgblog/goserver/pojo"
  "errors"
  "encoding/json"
)

func init(){
  ghttp.GetServer().BindHandler("post:/api/draft",     SaveDraft)
  ghttp.GetServer().BindHandler("get:/api/drafts",     GetAllDraft)
  ghttp.GetServer().BindHandler("patch:/api/draft/:aid",     UpdateDraft)
}


func SaveDraft(r *ghttp.Request){
    article := model.Article{}
    article.Content = r.PostForm.Get("content")
    article.Title = r.PostForm.Get("title")
    article.IsPublish = false
    article.Date = time.Now()
    article.CommentN = 0
    article.Tags = r.GetPostArray("tags")
    if saveArticle(article){
      util.SucceedRet("save ok",r)
    }else{
      util.FailRet(402,"save fail",r)
    }
}


func UpdateDraft(r *ghttp.Request){
  article := model.Article{}
  article.Content = r.PostForm.Get("content")
  article.Title = r.PostForm.Get("title")
  article.IsPublish = false
  article.Date = time.Now()
  article.CommentN = 0
  article.Tags = r.GetPostArray("tags")
  aidStr := r.Get("aid")
  if aidStr!=""{
    aidv,err := strconv.ParseInt(aidStr,10,32)
    if err!= nil {
      util.FailRet(402,"param error",r)
    }else{
      article.Aid = int32(aidv)
      if updateArticle(article){
        util.SucceedRet("update ok",r)
      }else{
        util.FailRet(402,"param error",r)
      }
    }
  }else{
    util.FailRet(402,"param error",r)
  }
}

func GetAllDraft(r *ghttp.Request){
  page := r.GetRequestMap()["page"]
  limit := r.GetRequestMap()["limit"]
  pageV,_ := strconv.Atoi(page)
  limitV,_ := strconv.Atoi(limit)
  if limitV==0 {
    limitV= 8
  }
  skip := limitV*(pageV-1)
  if data,err := getDraft(limitV,skip);err== nil {
    jsonData,_ := json.Marshal(data)
    util.SucceedRet(string(jsonData),r)
  }else{
    util.FailRet(402,"query draft fail",r)
  }
}


func getDraft(limit,skip int)([]pojo.ArticleResp,error){
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
  qryData := make([]model.Article,0)
  err := cu.Find(bson.M{"isPublish":false}).Sort("-date").Limit(limit).Skip(skip).All(&qryData)
  if err!= nil {
    return nil,err
  }else{
    retData := make([]pojo.ArticleResp,0)
    for idx,d := range qryData{
      r := pojo.ArticleResp{}
      if util.CopyProperty(&d,&r) {
        retData[idx] = r
      }else {
        return nil, errors.New("copy  fail")
      }
    }
    return retData,nil
  }
}
