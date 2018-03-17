package main  
//github  https://github.com/ssont/golang-open  
import (  
  "encoding/json"  
  "fmt"  
  "log"  
  "strconv"  
  "strings"  
  "time"  

  "github.com/astaxie/beego"  
  "github.com/astaxie/beego/session"  
)  

var globalSessions *session.Manager  

type Iindex struct {  
  beego.Controller  
}  

func init() {  
  config := fmt.Sprintf(`{"cookieName":"gosessionid","gclifetime":%d,"enableSetCookie":true}`, 3600*24) //  
  conf := new(session.ManagerConfig)  
  if err := json.Unmarshal([]byte(config), conf); err != nil {  
	  log.Fatal("json decode error", err)  
  }  
  globalSessions, _ = session.NewManager("memory", conf)  
  go globalSessions.GC()  
}  

func main() {  
  beego.BConfig.Listen.ServerTimeOut = 10 //设置 HTTP 的超时时间，默认是 0，不超时。  
	
  beego.BConfig.Listen.HTTPPort = 1000 //应用监听端口，默认为 8080。  

  beego.BConfig.AppName = "ich练习go Web编程"           //应用名称，默认是 beego。通过 bee new 创建的是创建的项目名。  
  beego.BConfig.ServerName = "QQ:1767311" //beego 服务器默认在请求的时候输出 server 为 beego。  

  beego.BConfig.WebConfig.Session.SessionName = "sessionID"         //存在客户端的 cookie 名称，默认值是 beegosessionID。  
  beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600 * 24  //session 过期时间，默认值是 3600 秒。  
  beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600 * 24 //session 默认存在客户端的 cookie 的时间，默认值是 3600 秒。  

  //beego.BConfig.WebConfig.Session.SessionDomain = "" //session cookie 存储域名, 默认空。  
  //beego.BConfig.WebConfig.ViewsPath = "admin" //模板路径，默认值是 views。  

  beego.Router("/*", &Iindex{}, "*:Count")  
  go beego.Run()  
  //=============================================  
  for { //死循环  
	  time.Sleep(10 * time.Second)  
  }  

}  

//网站访问计数器  
func (this *Iindex) Count() {  
  path_url := this.Ctx.Request.URL.String()  
  fmt.Println("get url:", path_url)  
  if path_url == "/favicon.ico" { //忽略此路由地址请求  
	  this.Ctx.WriteString("")  
	  this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html")  
	  return  
  }  

  //this.Ctx.Request  //这里面有大家所需要一切客户信息  
  fmt.Printf("===%v===\n", this.Ctx.Request)  

  Client_Host := this.Ctx.Request.Host                           //访问域名  
  Client_Method := this.Ctx.Request.Method                       //请求方式  
  Client_User_Agent := this.Ctx.Request.Header.Get("User-Agent") //请求头  
  Client_IP := this.Ctx.Request.Header.Get("Remote_addr")        //客户端IP  
  Client_Referer := this.Ctx.Request.Header.Get("Referer")       //来源  
  if len(Client_IP) <= 7 {  
	  Client_IP = this.Ctx.Request.RemoteAddr //获取客户端IP  
  }  
  if strings.Contains(Client_IP, ":") {  
	  ip_boolA, ip_dataA := For_IP(string(Client_IP)) //获取IP  
	  if ip_boolA {  
		  Client_IP = ip_dataA  
	  }  
  }  

  this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html")  
  this.Ctx.WriteString("golang网站流量统计  上</br>\n")  
  this.Ctx.WriteString("QQ:29295842</br>\n")  

  this.Ctx.WriteString(fmt.Sprintf("=====客户端IP:%v======</br>\n", Client_IP))  
  this.Ctx.WriteString(fmt.Sprintf("=====访问域名:%v======</br>\n", Client_Host))  
  this.Ctx.WriteString(fmt.Sprintf("=====请求路径:%v======</br>\n", path_url))  
  this.Ctx.WriteString(fmt.Sprintf("=====来源来路:%v======</br>\n", Client_Referer))  
  this.Ctx.WriteString(fmt.Sprintf("=====请求方式:%v======</br>\n", Client_Method))  
  this.Ctx.WriteString(fmt.Sprintf("=====请求头:%v======</br>\n", Client_User_Agent))  
  this.Ctx.WriteString(fmt.Sprintf("=====访问次数:%v======</br>\n", this.Cookie_session()))  
  //后面就是数据存贮  可以多种模式  

  return  
}  

func For_IP(valuex string) (bool, string) {  
  data_list := strings.Split(valuex, ":")  
  if len(data_list) >= 2 {  
	  return true, data_list[0]  
  }  
  return false, ""  
}  

func (this *Iindex) Cookie_session() int { //id统计  PV  这样统计只能针对单个浏览器有效  
  pv := 0  
  //=====================  
  //Cookie 统计法  
  cook := this.Ctx.GetCookie("countnum") //获取Cookie  
  if cook == "" {  
	  this.Ctx.SetCookie("countnum", "1", "/")  
	  pv = 1  
  } else {  
	  xx, err := strconv.Atoi(cook)  
	  if err == nil {  
		  pv = xx + 1  
		  this.Ctx.SetCookie("countnum", strconv.Itoa(pv), "/")  
	  } else {  
		  pv = 0  
	  }  
  }  
  return pv  
  //=====================  
  //session 统计法  
  sess, _ := globalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)  
  ct := sess.Get("countnum")  
  if ct == nil {  
	  sess.Set("countnum", 1)  
	  pv = 1  
  } else {  
	  pv = ct.(int) + 1  
	  sess.Set("countnum", pv)  
  }  
  return pv  
}  