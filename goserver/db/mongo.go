package db

import (
  "sync"
  "gopkg.in/mgo.v2"
)

const(
  DB_URL = "localhost:27017"
  DB_INFO = "my-blog"
  COMMENTS_C = "comments"
  ARTICLES_C = "articles"
  SEQUENCES_C = "sequences"
  USERS_C = "users"
)

func init(){

}

var once sync.Once
var dbsessionPool *sync.Pool

func newMongoSession()interface{}{
  sess,err := mgo.Dial(DB_URL)
  if err!= nil {
    return nil
  }else{
    return sess
  }
}

func GetDbPool() *sync.Pool{
  once.Do(func() {
    dbsessionPool = new(sync.Pool)
    dbsessionPool.New = newMongoSession
  })
  return dbsessionPool
}


func GetMongoConnFromPool()*mgo.Session{
  return GetDbPool().Get().(*mgo.Session)
}


func ReturnMongoConn(sess *mgo.Session){
  GetDbPool().Put(sess)
}

