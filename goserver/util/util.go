package util

import (
  "crypto/sha1"
  "encoding/hex"
  "time"
  "strings"
  "bytes"
  "gitee.com/johng/gf/g/net/ghttp"
  "github.com/magicgravity/vgblog/goserver/model"
  "github.com/magicgravity/vgblog/goserver/pojo"
  "reflect"
  "fmt"
  "strconv"
)

func Sha1(pwd ,salt string) string {
  sha := sha1.New()
  _,err :=sha.Write(([]byte)(pwd+salt))
  if err != nil {
    return ""
  }else {
    cs := sha.Sum(nil)
    return hex.EncodeToString(cs)
  }
}



func RandSalt()string{
  return ""
}

//2017-04-25T07:11:41.000Z
//YYYYMMDD HH:MI:SS
func ParseDate(datestr string)(time.Time,error){
  return time.Parse("20060102 15:04:05",datestr)
}


/*
格式化时间
 */
func FormatDate(date time.Time,fmtype string ) string{
  uperType := strings.ToUpper(fmtype)
  fmtStr := date.Format("20060102150405")

  buf :=bytes.Buffer{}
  if strings.Contains(uperType,"YYYY") {
    buf.WriteString(strings.Replace(uperType,"YYYY",fmtStr[0:4],-1))
  }else{
    buf.WriteString(uperType)
  }
  if strings.Contains(buf.String(),"MM") {
    month :=strings.Replace(buf.String(),"MM",fmtStr[4:6],-1)
    buf.Reset()
    buf.WriteString(month)
  }

  if strings.Contains(buf.String(),"DD") {
    day := strings.Replace(buf.String(),"DD",fmtStr[6:8],-1)
    buf.Reset()
    buf.WriteString(day)
  }

  if strings.Contains(buf.String(),"HH") {
    hour := strings.Replace(buf.String(),"HH",fmtStr[8:10],-1)
    buf.Reset()
    buf.WriteString(hour)
  }

  if strings.Contains(buf.String(),"MI") {
    minute := strings.Replace(buf.String(),"MI",fmtStr[10:12],-1)
    buf.Reset()
    buf.WriteString(minute)
  }

  if strings.Contains(buf.String(),"SS") {
    second := strings.Replace(buf.String(),"SS",fmtStr[12:],-1)
    buf.Reset()
    buf.WriteString(second)
  }

  return buf.String()
}

const (
  CommonMsg = "{\"msg\":\"some error happen\"}"
  SucceedMsg = "{\"msg\":\"succeed!\"}"
)



func FailRet(failCode int,failMsg string,r *ghttp.Request){
  if failCode>0 && failCode!=200 {
    r.Response.WriteHeader(failCode)
  }else{
    r.Response.WriteHeader(404)
  }
  if failMsg!=""{
    r.Response.WriteJson(failMsg)
  }else{
    r.Response.WriteJson(CommonMsg)
  }
}


func SucceedRet(msg string,r *ghttp.Request){
  r.Response.WriteHeader(200)
  if msg!= ""{
    r.Response.WriteJson(msg)
  }else{
    r.Response.WriteJson(SucceedMsg)
  }
}


func CopyArticleProp(src model.Article)pojo.ArticleResp{
  rr := pojo.ArticleResp{}
  rr.Id = src.Id_.Hex()
  rr.Content = src.Content
  rr.Title = src.Title
  rr.Aid = src.Aid
  rr.CommentN = src.CommentN
  rr.IsPublish = src.IsPublish
  copy(rr.Tags,src.Tags)
  rr.Date = FormatDate(src.Date,"YYYYMMDD HH:MI:SS")

  return rr
}


func BatchCopyArticleProp(src []model.Article)[]pojo.ArticleResp{
  ret := make([]pojo.ArticleResp,len(src))
  for idx,v :=range src {
    ret[idx] = CopyArticleProp(v)
  }

  return ret
}


func CopyProperty(src ,des interface{})bool{
  if src==nil || des ==nil {
    return false
  }
  defer func(){
    if r := recover();r!= nil {
      fmt.Errorf("copy property fail ! ,reason maybe %v \r\n",r)
    }
  }()

  srcType := reflect.TypeOf(src)
  desType := reflect.TypeOf(des)

  //srcType Kind: struct
  //fmt.Println("srcType Kind:",srcType.Kind().String())
  ////srcType Name: User
  //fmt.Println("srcType Name:",srcType.Name())
  ////fmt.Println("srcType Elem:",srcType.Elem().String())
  //fmt.Println("srcType NumField:",srcType.NumField())
  //fmt.Println("srcType Field 2 Name:",srcType.Field(2).Name)
  //fmt.Println("srcType Field 2 Type:",srcType.Field(2).Type)

  //拷贝目标对象必须是指针类型
  if desType.Kind().String()!="ptr" {
    return false
  }
  //如果不是指针也不是结构体
  if srcType.Kind().String()!="ptr" && srcType.Kind().String()!="struct"{
    return false
  }
  var srcFieldNum int
  if srcType.Kind()==reflect.Ptr{
    srcFieldNum = srcType.Elem().NumField()
  }else{
    srcFieldNum = srcType.NumField()
  }


  desVal := reflect.ValueOf(des).Elem()
  srcVal := reflect.ValueOf(src).Elem()
  //fmt.Printf("srcFieldNum >> %d \r\n" ,srcFieldNum)
  for i:=0;i<srcFieldNum;i++ {
    var ftype reflect.Type
    var fname string
    if srcType.Kind()==reflect.Ptr{
      ftype = srcType.Elem().Field(i).Type
      fname = srcType.Elem().Field(i).Name
    }else{
      ftype = srcType.Field(i).Type
      fname = srcType.Field(i).Name
    }

    //fmt.Printf("src ftype == > %v,fname ==> %v \r\n",ftype,fname)

    if dfield,ok := desType.Elem().FieldByName(fname);ok && desVal.FieldByName(fname).CanSet() {

      dtype := dfield.Type
      //fmt.Printf("ftype >> %v ,fname >> %v , dtype >> %v \r\n" ,ftype,fname,dtype)
      if dtype.String()==ftype.String() {
        //直接复制
        desVal.FieldByName(fname).Set(srcVal.FieldByName(fname))
      }else{
        fmt.Printf("dtype >> %v,ftype >> %v \r\n",dtype.String(),ftype.String())
        if dtype.String()=="string" && ftype.String()=="time"{
          //time 转 string
          ftime := srcVal.FieldByName(fname).Interface().(time.Time)
          ftimeStr := FormatDate(ftime,"YYYYMMDD HHMISS")
          desVal.FieldByName(fname).SetString(ftimeStr)
        }else if dtype.String()=="string" && (
          ftype.String()=="int"  ){
          //int 转 string
          fint := srcVal.FieldByName(fname).Interface().(int)
          fstr := strconv.Itoa(fint)
          desVal.FieldByName(fname).SetString(fstr)
        }else if dtype.String()=="string" && (
          ftype.String()=="int8" || ftype.String()=="int16"  || ftype.String()=="int32" ||  ftype.String()=="int64" ){
          //int 转 string
          var fstr string
          switch ftype.Kind() {
            case reflect.Int8:
              fint := srcVal.FieldByName(fname).Interface().(int8)
              fstr = strconv.FormatInt(int64(fint),10)
            case reflect.Int16:
              fint := srcVal.FieldByName(fname).Interface().(int16)
              fstr = strconv.FormatInt(int64(fint),10)
            case reflect.Int32:
              fint := srcVal.FieldByName(fname).Interface().(int32)
              fstr = strconv.FormatInt(int64(fint),10)
            case reflect.Int64:
              fint := srcVal.FieldByName(fname).Interface().(int64)
              fstr = strconv.FormatInt(fint,10)
            default:
              continue
          }

          desVal.FieldByName(fname).SetString(fstr)
        }else if dtype.String()=="string" &&(
          ftype.String()=="uint" || ftype.String()=="uint8" || ftype.String()=="uint16"  || ftype.String()=="uint32" || ftype.String()=="uint64"){
          //uint 转 string
          //fmt.Printf("~~~~0     %v \r\n",srcVal.FieldByName(fname).Interface())
          var fstr string
          switch ftype.String() {
            case "uint":
              fint := srcVal.FieldByName(fname).Interface().(uint)
              fstr = strconv.FormatUint(uint64(fint),10)
            case "uint8":
              fint := srcVal.FieldByName(fname).Interface().(uint8)
              fstr = strconv.FormatUint(uint64(fint),10)
            case "uint16":
              fint := srcVal.FieldByName(fname).Interface().(uint16)
              fstr = strconv.FormatUint(uint64(fint),10)
            case "uint32":
              fint := srcVal.FieldByName(fname).Interface().(uint32)
              fstr = strconv.FormatUint(uint64(fint),10)
            case "uint64":
              fint := srcVal.FieldByName(fname).Interface().(uint64)
              fstr = strconv.FormatUint(fint,10)
            default:
              continue
          }

          desVal.FieldByName(fname).SetString(fstr)
        }else if dtype.String()=="string" && ftype.String()=="bool" {
          srcBool :=srcVal.FieldByName(fname).Interface().(bool)
          desVal.FieldByName(fname).SetString(strconv.FormatBool(srcBool))

        }else if dtype.String()=="string" && (ftype.Kind()==reflect.Float32 || ftype.Kind()==reflect.Float64) {
          var srcStr string
          switch ftype.Kind() {
            case reflect.Float64:
              srcF := srcVal.FieldByName(fname).Interface().(float64)
              srcStr = strconv.FormatFloat(srcF,'f',5,64)
            case reflect.Float32:
              srcF := srcVal.FieldByName(fname).Interface().(float32)
              srcStr = strconv.FormatFloat(float64(srcF),'f',5,32)
            default:
              continue
          }
          desVal.FieldByName(fname).SetString(srcStr)

        }
      }
    }else{
      continue
    }
  }

  return true
  //srcVal := reflect.ValueOf(src)
  //desVal := reflect.ValueOf(des)
}
