package model

import (
  "gopkg.in/mgo.v2/bson"
  "time"
)

type User struct {
  Id_             bson.ObjectId `bson:"_id"`
  Name string   `bson:"name"`
  Password string `bson:"password"`
  Salt string `bson:"salt"`
}


type Article struct {
  Id_             bson.ObjectId `bson:"_id"`
  Aid             int32  `bson:"aid"`
  Title           string `bson:"title"`
  Content         string `bson:"content"`
  Date            time.Time `bson:"date"`
  IsPublish       bool `bson:"isPublish"`
  CommentN        int32 `bson:"comment_n"`
  Tags            []string `bson:"tags"`
}

type Comment struct {
  Id_             bson.ObjectId `bson:"_id"`
  ImgName         string    `bson:"imgName"`
  Name            string    `bson:"name"`
  Address         string    `bson:"address"`
  Date            time.Time `bson:"date"`
  Content         string    `bson:"content"`
  ArticleId       int32     `bson:"articleId"`
  Like            int32     `bson:"like"`
}
