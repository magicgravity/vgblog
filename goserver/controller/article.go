package controller

import (
  "gitee.com/johng/gf/g/net/ghttp"
  "github.com/magicgravity/vgblog/goserver/model"
  "github.com/magicgravity/vgblog/goserver/db"
  "time"
  "gopkg.in/mgo.v2/bson"
  "math"
  "gopkg.in/mgo.v2"
  "github.com/magicgravity/vgblog/goserver/pojo"
  "errors"
  "strconv"
  "github.com/magicgravity/vgblog/goserver/util"
  "encoding/json"
  "strings"
  "github.com/golang/glog"
)

func init(){
  ghttp.GetServer().BindHandler("post:/api/article",     PublishArticle)
  ghttp.GetServer().BindHandler("get:/api/article/:aid",     DetailArticle)
  ghttp.GetServer().BindHandler("delete:/api/article/:aid",     DeleteArticle)
  ghttp.GetServer().BindHandler("patch:/api/article/:aid",     UpdateArticle)
  ghttp.GetServer().BindHandler("get:/api/articles",     GetArticles)
  ghttp.GetServer().BindHandler("get:/api/someArticles",     SearchArticles)
}


func PublishArticle(r *ghttp.Request){
  article := model.Article{}
  article.Content = r.PostForm.Get("content")
  article.Title = r.PostForm.Get("title")
  article.IsPublish = true
  article.Date = time.Now()
  article.CommentN = 0
  if curMax,err := findMaxAid();err== nil {
    if curMax <math.MaxInt32-1 {
      article.Aid = (curMax + 1)
    }else{
      panic("article aid reach max value!")
    }
  }else{
    util.FailRet(0,"",r)
  }
  if saveArticle(article) {
    util.SucceedRet("",r)
  }else{
    util.FailRet(0,"",r)
  }
}

func findMaxAid()(int32,error){
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
  article := model.Article{}
  err := cu.Find(bson.M{}).Sort("-aid").Limit(1).One(&article)
  if err!= nil {
    return 0,err
  }else{
    return article.Aid,nil
  }
}

func saveArticle(article model.Article)bool{
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
  err := cu.Insert(article)
  if err!= nil {
    return false
  }else{
    ensureIndex(cu)
    return true
  }
}

func ensureIndex(cu *mgo.Collection){
  index := mgo.Index{
             Key: []string{"aid"},
             Unique: true,
             //DropDups: true,
             //Background: true, // See notes.
             //Sparse: true,
         }
  cu.EnsureIndex(index)
}

func DetailArticle(r *ghttp.Request){
  ret,err := findOneArticle(r.Get("aid"))
  if err==nil {
    retData,_ := json.Marshal(ret)
    util.SucceedRet(string(retData),r)
  }else{
    util.FailRet(0,"",r)
  }
}

func findOneArticle(aid string) (pojo.ArticleResp,error){
  if aid==""{
    return pojo.ArticleResp{},errors.New("not found")
  }else{
    aidV,_ := strconv.ParseInt(aid,10,32)
    sess := db.GetMongoConnFromPool()
    defer db.ReturnMongoConn(sess)
    cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
    a := model.Article{}
    err := cu.Find(bson.M{"aid":aidV}).Limit(1).One(&a)
    if err!= nil {
      return pojo.ArticleResp{},errors.New("not found")
    }else{
      rr := pojo.ArticleResp{}
      rr.Id = a.Id_.Hex()
      rr.Content = a.Content
      rr.Title = a.Title
      rr.Aid = a.Aid
      rr.CommentN = a.CommentN
      rr.IsPublish = a.IsPublish
      copy(rr.Tags,a.Tags)
      rr.Date = util.FormatDate(a.Date,"YYYYMMDD HH:MI:SS")
      return rr,nil
    }

  }
}

func DeleteArticle(r *ghttp.Request){
  if deleteArticleByAid(r.Get("aid")) {
    r.Response.WriteHeader(200)
  }else{
    r.Response.WriteHeader(404)
  }
}

func deleteArticleByAid(aid string)bool{
  if aid==""{
    return false
  }else{
    aidV,err := strconv.ParseInt(aid,10,32)
    if err!= nil {
      return false
    }
    sess := db.GetMongoConnFromPool()
    defer db.ReturnMongoConn(sess)
    cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
    err = cu.Remove(bson.M{"aid":aidV})
    if err!= nil {
      return false
    }else{
      return true
    }
  }
}

func UpdateArticle(r *ghttp.Request){
  aid := r.Get("aid")
  aidV,err := strconv.ParseInt(aid,10,32)
  if err!= nil {
    util.FailRet(0,err.Error(),r)
    return
  }else{
    article := model.Article{}
    article.Aid = int32(aidV)
    article.Title = r.Form.Get("title")
    article.Content = r.Form.Get("content")
    article.IsPublish = true
    article.Date = time.Now()
    article.Tags = r.GetRequestArray("tags")
    if updateArticle(article) {
      util.SucceedRet("",r)
    }else{
      util.FailRet(0,"",r)
    }
  }
}

func updateArticle(article model.Article)bool{
  if article.Aid<1 {
    return false
  }else{
    sess := db.GetMongoConnFromPool()
    defer db.ReturnMongoConn(sess)
    cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
    err := cu.Update(bson.M{"aid":article.Aid},bson.M{
      "title":article.Title,
      "content":article.Content,
      "date":article.Date,
      "isPublish":article.IsPublish,
      "tags":article.Tags,
    })
    if err!= nil{
      return false
    }else{
      return true
    }
  }
}

func GetArticles(r *ghttp.Request){
  page := r.GetRequestMap()["page"]
  value := r.GetRequestMap()["value"]
  limit := r.GetRequestMap()["limit"]
  limitV,_ := strconv.Atoi(limit)
  if limitV <0  {
    limitV = 4
  }
  pageV,_ := strconv.Atoi(page)
  skip := limitV*(pageV-1)

  tags := []string{value}
  if ret := getArticles(tags,limitV,skip,"全部"==value);ret!= nil {
    retData,_ :=  json.Marshal(ret)
    util.SucceedRet(string(retData),r)
  }else{
    util.FailRet(0,"",r)
  }
}

func getArticles(tags []string,limit,skip int,allFlag bool)[]pojo.ArticleResp{
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
  articles := make([]model.Article,0)
  var err error
  if allFlag{
    //查全部
    err = cu.Find(bson.M{"isPublish":true}).Sort("-date").Limit(limit).Skip(skip).All(&articles)
  }else{
    err = cu.Find(bson.M{"isPublish":true,"tags":tags}).Sort("-date").Limit(limit).Skip(skip).All(&articles)
  }
  if err!= nil {
    glog.Errorf("getArticles fail ,reason ==> %v \r\n",err)
    return nil
  }else{
    return util.BatchCopyArticleProp(articles)
  }
}

func SearchArticles(r *ghttp.Request){
  key := r.GetRequestMap()["key"]
  page := r.GetRequestMap()["page"]
  value := r.GetRequestMap()["value"]
  pageV,_ := strconv.Atoi(page)
  if pageV<1{
    pageV = 1
  }
  skip := 4*(pageV-1)

  if ret,err := searchArticles(key,value,skip);err== nil {
    retData,_ :=json.Marshal(ret)
    util.SucceedRet(string(retData),r)
  }else{
    util.FailRet(0,err.Error(),r)
  }
}


func searchArticles(key,value string,skip int)([]pojo.ArticleResp,error){
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)

  qryRets := make([]model.Article,0)
  var err error
  switch key {
  case "tags":
    tagsArr := strings.Split(value," ")
    err = cu.Find(bson.M{"tags":bson.M{"$all":tagsArr}}).Sort("-date").Limit(4).Skip(skip).All(&qryRets)
  case "title":
    err = cu.Find(bson.M{"title":
            bson.M{"$regex":value,
                   "$options":"$i"},
                   "isPublish":true}).Sort("-date").Limit(4).Skip(skip).All(&qryRets)
  case "date":
    beginDate,err := util.ParseDate(value)
    if err== nil {
      endDate := beginDate.Add(24*time.Hour)
      err = cu.Find(bson.M{"date":
        bson.M{"$gte":beginDate,
               "$lt":endDate}}).Sort("-date").Limit(4).Skip(skip).All(&qryRets)
    }
  default:
    err = errors.New("key type is not right")
  }
  if err==nil {
    return util.BatchCopyArticleProp(qryRets),nil
  }else{
    return nil,err
  }
}

