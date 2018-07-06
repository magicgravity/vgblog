package controller

import (
  "gitee.com/johng/gf/g/net/ghttp"
  "github.com/magicgravity/vgblog/goserver/db"
  "gopkg.in/mgo.v2/bson"
  "github.com/magicgravity/vgblog/goserver/util"
  "github.com/magicgravity/vgblog/goserver/model"
  "time"
  "strconv"
  "github.com/magicgravity/vgblog/goserver/pojo"
  "encoding/json"
)

func init(){
  ghttp.GetServer().BindHandler("post:/api/comment",     PostComment)
  ghttp.GetServer().BindHandler("get:/api/comments",     GetComments)
  ghttp.GetServer().BindHandler("patch:/api/comments/:id",     UpdateComment)
}

func PostComment(r *ghttp.Request){
  name := r.PostForm.Get("name")
  articleId := r.PostForm.Get("articleId")
  address := r.PostForm.Get("address")
  if ok,_ :=findCommentByNameArticleId(name,articleId,address) ; ok{
    comment := model.Comment{}
    comment.Address = address
    comment.Date = time.Now()
    comment.Like = 0
    comment.Name = name
    comment.Content = r.PostForm.Get("content")

    if err := postComment(comment);err==nil {
      //send mail
      updateArticleComment(articleId)

      util.SucceedRet("save comment successful",r)
    }else{
      util.FailRet(401,"save comment fail",r)
    }
  }else{
    util.FailRet(403,"用户名已经存在",r)
  }
}

func findCommentByNameArticleId(name,articleId,address string)(bool,model.Comment){
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.COMMENTS_C)
  ret := model.Comment{}
  err := cu.Find(bson.M{"name":name,"articleId":articleId}).One(&ret)
  if err!=nil || ret.Address=="" || ret.Address!=address {
    return false,ret
  }else{
    return true,ret
  }
}

func postComment(comment model.Comment)error{
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.COMMENTS_C)
  err := cu.Insert(comment)
  if err!= nil {
    return err
  }else{
    cu2 := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
    err := cu2.Update(bson.M{"aid":comment.ArticleId},bson.M{"$inc":bson.M{"comment_n":1}})
    return err
  }
}



func GetComments(r *ghttp.Request){
  articleId := r.PostForm.Get("id")
  sort := r.PostForm.Get("sort")
  if data,err := getComment(articleId,sort);err==nil {
    dataStr,_ := json.Marshal(data)
    util.SucceedRet(string(dataStr),r)
  }else{
    util.FailRet(401,err.Error(),r)
  }
}

func getComment(id ,sort string)([]pojo.CommentResp,error){
  articleId,err := strconv.ParseInt(id,10,32)
  if err!= nil {
    return nil,err
  }else{
    qryRets := make([]model.Comment,0)
    var err error
    sess := db.GetMongoConnFromPool()
    defer db.ReturnMongoConn(sess)
    cu := sess.DB(db.DB_INFO).C(db.COMMENTS_C)
    switch sort {
      case "date":
        err = cu.Find(bson.M{"articleId":articleId}).Select(bson.M{"name":1,
                                                             "date":1,
                                                             "content":1,
                                                             "like":1,
                                                             "imgName":1}).Sort("-date").All(&qryRets)
      case "like":
        err = cu.Find(bson.M{"articleId":articleId}).Select(bson.M{"name":1,
                                                                   "date":1,
                                                                   "content":1,
                                                                   "like":1,
                                                                   "imgName":1}).Sort("-like").All(&qryRets)
      default:
        err = cu.Find(bson.M{"articleId":articleId}).Select(bson.M{"name":1,
                                                                   "date":1,
                                                                   "content":1,
                                                                   "like":1,
                                                                   "imgName":1}).All(&qryRets)

    }
    if err!= nil {
      return nil,err
    }else{
      retComms := make([]pojo.CommentResp,len(qryRets))
      for idx,q :=range qryRets{
        ctmp := pojo.CommentResp{}
        util.CopyProperty(&q,&ctmp)
        retComms[idx] = ctmp
      }
      return retComms,nil
    }
  }
}


func UpdateComment(r *ghttp.Request){
  option := r.PostForm.Get("option")
  id := r.Get("id")
  if option=="add" || option=="drop" {
    if updateComment(id, option == "add"){
      util.SucceedRet("succeed in updating like",r)
    }else{
      util.FailRet(401,"fail in updating like",r)
    }
  }else {
    util.FailRet(401, "param error", r)
  }
}


func updateComment(cid string,uflag bool)bool{
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.COMMENTS_C)
  var err error
  if uflag {
    err = cu.Update(bson.M{"_id": bson.ObjectIdHex(cid)}, bson.M{"$inc": bson.M{"like": 1}})
  }else{
    err = cu.Update(bson.M{"_id": bson.ObjectIdHex(cid)}, bson.M{"$inc": bson.M{"like": -1}})
  }
  if err==nil {
    return true
  }else{
    return false
  }
}


func updateArticleComment(aid string)bool{
  if aid==""{
    return false
  }else {
    aidV, err := strconv.ParseInt(aid, 10, 32)
    if err!= nil {
      return false
    }
    sess := db.GetMongoConnFromPool()
    defer db.ReturnMongoConn(sess)
    cu := sess.DB(db.DB_INFO).C(db.ARTICLES_C)
    err = cu.Update(bson.M{"aid":aidV},bson.M{"$inc":bson.M{"comment_n":1}})
    if err== nil {
      return true
    }else{
      return false
    }
  }
}
