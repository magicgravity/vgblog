package pojo


type LoginResp struct {
  Id string   `json:"id"`
  Name string `json:"name"`
  Token string  `json:"token"`
  RetCode string `json:"retCode"`
}


type ArticleResp struct {
  Id string   `json:"id"`
  Aid int32   `json:"aid"`
  Title string  `json:"title"`
  Content string  `json:"content"`
  Date string `json:"date"`
  IsPublish bool `json:"isPublish"`
  CommentN int32  `json:"commentN"`
  Tags []string   `json:"tags"`
}


type CommentResp struct {
  Id string   `json:"id"`
  ImgName string `json:"imgName"`
  Name string   `json:"name"`
  Address string  `json:"address"`
  Date string   `json:"date"`
  Content string  `json:"content"`
  ArticleId int32 `json:"articleId"`
  Like  int32 `json:"like"`
}
