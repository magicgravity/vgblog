package controller


import (
  "gitee.com/johng/gf/g/net/ghttp"
  "gopkg.in/mgo.v2/bson"
  "github.com/magicgravity/vgblog/goserver/db"
  "github.com/magicgravity/vgblog/goserver/model"
  "github.com/golang/glog"
  "github.com/magicgravity/vgblog/goserver/util"
  "github.com/magicgravity/vgblog/goserver/pojo"
  "encoding/json"
)

func init(){
  ghttp.GetServer().BindHandler("post:/api/login",     Login)
  ghttp.GetServer().BindHandler("post:/api/user",     UpdateUser)
}



func Login(r *ghttp.Request){
  name := r.PostForm.Get("name")
  pwd := r.PostForm.Get("password")

  resp := login(name,pwd)
  if resp.RetCode=="ok" {
    respData,err := json.Marshal(resp)
    if err== nil {
      r.Response.WriteJson(string(respData))
    }else{
      r.Response.WriteHeader(500)
    }
  }else{
    r.Response.WriteHeader(401)
  }
}

func login(name ,pwd string) *pojo.LoginResp{
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.USERS_C)

  u := new(model.User)
  err := cu.Find(bson.M{"name": name}).One(u)
  if err != nil{
    glog.Error("find user by username fail ",err)
    resp := new(pojo.LoginResp)
    resp.RetCode = err.Error()
    return resp
  }else{
    if u.Password == util.Sha1(pwd,u.Salt) {

      /*
      token 生成
      jwt.sign({
        id: id,
        name: name
      }, secret.cert, { expiresIn: '7d' })
       */

      resp := new(pojo.LoginResp)
      resp.Id = u.Id_.Hex()
      resp.Name = u.Name
      resp.RetCode = "ok"
      resp.Token = ""   //todo
      return resp
    }else{
      resp := new(pojo.LoginResp)
      resp.RetCode = "401"
      return resp
    }
  }
}


func UpdateUser(r *ghttp.Request){
  salt := util.RandSalt()
  user := new(model.User)
  user.Id_ = bson.ObjectIdHex(r.PostForm.Get("id"))
  user.Name = r.PostForm.Get("name")
  user.Salt = salt
  user.Password = util.Sha1(r.PostForm.Get("password"),salt)

  if updateUser(user) {
    r.Response.WriteHeader(200)
  }else{
    r.Response.WriteHeader(401)
  }
}


func updateUser(u *model.User)bool {
  sess := db.GetMongoConnFromPool()
  defer db.ReturnMongoConn(sess)
  cu := sess.DB(db.DB_INFO).C(db.USERS_C)

  err := cu.Update(bson.M{"_id":u.Id_},bson.M{"$set":bson.M{
    "password":u.Password,
    "salt":u.Salt,
  }})
  if err!= nil {
    glog.Error("update user info fail >> ",err)
    return false
  }else{
    return true
  }
}
