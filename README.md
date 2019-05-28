# Github_Hunter_By_Golang
最新版的Github Hunter由Golang重写，加入了新的一些小特性。
本工具主要是搜索Github存在的敏感信息，如用户名、密码和基础设置配置信息等。
## 运行要求
Golang 1.11.5 <br>
## 系统支持
Linux,MacOS,Windows<br>
## 安装使用方法
1.`git clone https://github.com/Hell0W0rld0/Github_Hunter_By_Golang.git`<br>
2.`cd Github_Hunter_By_Golang`<br>
3.`如果你没有Golang的运行环境，需要从Golang官方网站下载，具体下载方法请自行搜索，谢谢！`<br>
4.`安装好Golang环境后，可以直接使用go run Github_Hunter.go 运行，也可以编译成可执行文件。`<br>
5.`运行或编译前需要使用go get来安装依赖`<br>
6.`go get github.com/gocolly/colly github.com/mattn/go-sqlite3 gopkg.in/cheggaaa/pb.v2 gopkg.in/gomail.v2 gopkg.in/ini.v1 `<br>
  `安装完成后即可使用 go build Github_Hunter.go 来生成可执行文件。`
## 设置
使用前，需要配置ini文件，把info.ini.example更改为info.ini,然后根据需要填写以下内容
### 例子
`[KEYWORD]`<br>
`keyword1 = 主关键词，如域名`<br>
`keyword2 = 主关键词`<br>
`keyword3 = 主关键词`<br>
`...etc`<br>
<br>
`[EMAIL]`<br>
`host = 邮件服务器`<br>
`user = 邮件用户名`<br>
`password = 邮件密码`<br>
<br>
`[SENDER]`<br>
`sender = 发送者的邮件地址`<br>
<br>
`[RECEIVER]`<br>
`receiver1 = 接收者邮件地址-1`<br>
`receiver2 = 接收者邮件地址-2`<br>
<br>
`[Github]`<br>
`user = Github 用户名`<br>
`password = Github 密码`<br>
<br>
`[PAYLOADS]`<br>
`p1 = Payload 1`<br>
`p2 = Payload 2`<br>
`p3 = Payload 3`<br>
`p4 = Payload 4`<br>
`p5 = Payload 5`<br>
`p6 = Payload 6`<br>
### 关键词 和 载荷设置
尽量设置2-5个主关键词，如baidu.com,BaiDu.com,BAIDU.com,baidu.COM,BAIDU.COM等。<br>
Payloads处，加入敏感关键词，如password、Password、Username、UserName、Database、Mysql等。<br>
## 运行
如已编译成可执行程序则，使用./Github_Hunter运行，Windows: Github_Hunter.exe。<br>
未编译则使用go run Github_Hunter.go运行。<br>
运行后搜索到敏感信息会发送邮件进行通知，如果长期监控，请使用计划任务crontab来设置运行频度。<br>
