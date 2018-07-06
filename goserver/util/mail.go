package util

import (
  "net/smtp"
  "strings"
  "bytes"
  "html/template"
)

const (
  SMTP_163_ADDR = "smtp.163.com:25"
  SENDER_MAIL = "blog_admin666@163.com"
  SENDER_MAIL_PASSWORD = "123456"
)

var common_mail_template = `div style="width: 90%; border: 2px solid lightgreen; margin: 1rem auto; padding: 1rem; text-align: center;">
        <p style="border-bottom: 1px dashed lightgreen; margin: 0;padding-bottom: 1rem; color: lightgreen; font-size: 1.25rem;">{{.title}}</p>
        <p style="margin: 1rem 0 0;">hello,{{.name}} &#x1f608</p>
        <p style="margin: 0 0 1rem;">{{.otherName}}{{.message}}</p>
        <p style="width: 70%; border-left: 4px solid lightgreen; padding: 1rem; margin: 0 auto 2rem; text-align: left;white-space: pre-line;">{{.content}}</p>
    <a href= {{.url}} style="text-decoration: none; background: lightgreen;color: #fff; height: 2rem; line-height: 2rem; padding: 0 1rem; display: inline-block; border-radius: 0.2rem;">前往查看</a>
        </div>`

func MakeMailContent(title,name,otherName,message,content,url string)string{
  tpl,err := template.New("mailTemplate").Parse(common_mail_template)
  if err!= nil {
    return ""
  }else{
    data := struct{
      title string
      name string
      otherName string
      message string
      content string
      url string
    }{
      title:title,
      name:name,
      otherName:otherName,
      message:message,
      content:content,
      url:url,
    }

    buffer := bytes.NewBuffer(nil)
    err := tpl.Execute(buffer,data)
    if err!= nil {
      return ""
    }else{
      return buffer.String()
    }

  }
}

func SendMail(recMail,subject,content string)bool{

  auth := smtp.PlainAuth("",SENDER_MAIL,SENDER_MAIL_PASSWORD,
    strings.Split(SMTP_163_ADDR,":")[0])

  buf := bytes.Buffer{}
  buf.WriteString("To: ")
  buf.WriteString(recMail)
  buf.WriteString(" \r\nFrom: ")
  buf.WriteString(SENDER_MAIL)
  buf.WriteString("\r\nSubject: ")
  buf.WriteString(subject)
  buf.WriteString("\r\nContent-Type:text/html;charset=UTF-8\r\n\r\n")
  buf.WriteString(content)
  //client.Text.W.WriteString(buf.String())
  //client.Text.W.Flush()
  err := smtp.SendMail(SMTP_163_ADDR,auth,SENDER_MAIL,[]string{recMail},buf.Bytes())
  if err!= nil {
    return false
  }else{
    return true
  }

  //client,err := smtp.Dial(SMTP_163_ADDR)
  //if err!= nil {
  //
  //}else{
  //  if err := client.Mail(SENDER_MAIL); err != nil {
  //
  //  }
  //  if err := client.Rcpt(recMail); err != nil {
  //
  //  }
  //
  //
  //}

}
