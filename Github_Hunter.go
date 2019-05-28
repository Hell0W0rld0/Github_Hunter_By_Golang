// Github_Hunter
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/cheggaaa/pb.v2"
	"gopkg.in/gomail.v2"
	"gopkg.in/ini.v1"
)

func main() {
	/*fmt.Println(`
	     #####                                  #     #
	    #     # # ##### #    # #    # #####     #     # #    # #    # ##### ###### #####
	    #       #   #   #    # #    # #    #    #     # #    # ##   #   #   #      #    #
	    #  #### #   #   ###### #    # #####     ####### #    # # #  #   #   #####  #    #
	    #     # #   #   #    # #    # #    #    #     # #    # #  # #   #   #      #####
	    #     # #   #   #    # #    # #    #    #     # #    # #   ##   #   #      #   #
	     #####  #   #   #    #  ####  #####     #     #  ####  #    #   #   ###### #    #    V2.1 Created by Allen
		`)
	*/
	banner()
	var keywords []string
	var warningList []string
	start := time.Now()
	f, err := ini.Load("info.ini")
	if err != nil {
		fmt.Println("配置文件读取失败或不存在文件！")
		os.Exit(0)
	}
	gUser := f.Section("Github").Key("user").String()
	gPassword := f.Section("Github").Key("password").String()
	receiver := f.Section("RECEIVER").Key("receiver").Strings(",")
	host := f.Section("EMAIL").Key("host").String()
	user := f.Section("EMAIL").Key("user").String()
	password := f.Section("EMAIL").Key("password").String()
	sender := f.Section("SENDER").Key("sender").String()
	page := f.Section("PAGE").Key("page").String()
	number, err := strconv.Atoi(page)
	c1 := hunterLogin(gUser, gPassword)
	keyKeywords := f.Section("KEYWORD").Keys()
	keyPayloads := f.Section("PAYLOADS").Keys()
	for _, keyKeyword := range keyKeywords {
		for _, keyPayload := range keyPayloads {
			keyword := keyKeyword.Value() + "+" + keyPayload.Value()
			keywords = append(keywords, keyword)
		}
	}
	links, codes := hunterGetData(c1, keywords, number)
	if _, err := os.Stat("hunter.db"); err == nil {
		fmt.Println("存在数据文件，开始进行新增数据查找......")
		for _, k := range keywords {
			searchPayload := strings.Split(k, "+")
			for i := 0; i < len(links); i++ {
				if (strings.Contains(codes[i], searchPayload[0])) && (strings.Contains(codes[i], searchPayload[1])) {
					codes[i] = strings.Replace(codes[i], searchPayload[0], `<em style="color:red">`+searchPayload[0]+"</em>", -1)
					codes[i] = strings.Replace(codes[i], searchPayload[1], `<em style="color:red">`+searchPayload[1]+"</em>", -1)
					if searchedURL, _, _ := compareDBURL(links[i]); searchedURL == "" {
						warningList = append(warningList, "<br><br><br>链接："+links[i]+"<br><br>")
						warningList = append(warningList, "发现的关键词："+`<em style="color:red">`+searchPayload[0]+"</em> and "+`<em style="color:red">`+searchPayload[1]+"</em>"+`<br>简要代码如下：<br><div style="border:1px solid #bfd1eb;background:#f3faff">`+codes[i]+"</div>")
						insertDB(links[i], codes[i])
					}
				}
			}
		}
		fmt.Println("操作完毕！")
	} else {
		fmt.Println("未发现数据库文件，创建并建立基线......")
		for _, k := range keywords {
			searchPayload := strings.Split(k, "+")
			for i := 0; i < len(links); i++ {
				if (strings.Contains(codes[i], searchPayload[0])) && (strings.Contains(codes[i], searchPayload[1])) {
					codes[i] = strings.Replace(codes[i], searchPayload[0], `<em style="color:red">`+searchPayload[0]+"</em>", -1)
					codes[i] = strings.Replace(codes[i], searchPayload[1], `<em style="color:red">`+searchPayload[1]+"</em>", -1)
					insertDB(links[i], codes[i])
				}
			}
		}
		initURLs, initCodes := initInfomation()
		for _, key := range keywords {
			searchPayload := strings.Split(key, "+")
			for j := 0; j < len(initURLs); j++ {
				if (strings.Contains(initCodes[j], searchPayload[0])) && (strings.Contains(initCodes[j], searchPayload[1])) {
					warningList = append(warningList, "<br><br><br>链接："+initURLs[j]+"<br><br>")
					warningList = append(warningList, "发现的关键词："+`<em style="color:red">`+searchPayload[0]+"</em> and "+`<em style="color:red">`+searchPayload[1]+"</em>"+`<br>简要代码如下：<br><div style="border:1px solid #bfd1eb;background:#f3faff">`+initCodes[j]+"</div>")
				}
			}
		}

		fmt.Println("操作完毕！")
	}
	if len(warningList) > 0 {
		message := strings.Join(warningList, "")
		number := strconv.Itoa(len(warningList) / 2)
		result := fmt.Sprintf(`Dear all<br><br>发现疑似信息泄露! 一共发现 <em style="color:red">%s</em> 条！`+message, number)
		sendMail(receiver, sender, result, host, user, password)
	} else {
		info := "Dear all<br><br>未发现任何新增敏感信息！"
		sendMail(receiver, sender, info, host, user, password)
	}
	end := time.Now()
	total := end.Sub(start)
	fmt.Printf("本次耗时：%.2f 秒", total.Seconds())
}

func hunterLogin(gUser, gPassword string) *colly.Collector {
	//获取Token
	var token string
	c := colly.NewCollector()
	c.OnHTML("input[name=authenticity_token]", func(e *colly.HTMLElement) {
		token = e.Attr("value")
	})
	c.Visit("https://github.com/login")

	//登陆Github
	_loginData := map[string]string{"commit": "Sign in", "utf8": "✓", "authenticity_token": token, "login": gUser, "password": gPassword}
	err := c.Post("https://github.com/session", _loginData)
	if err != nil {
		fmt.Println("登陆出现问题：", err)
		panic(err)
	}
	c.Visit("https://github.com")

	c.OnResponse(func(f *colly.Response) {
		if strings.Contains(string(f.Body), `"Sign in, switch to sign up"`) {
			fmt.Println("登陆失败，请检查用户名和密码！")
			os.Exit(0)
		} else {
			fmt.Println("登陆成功，开始收集信息.......")

		}
	})
	c.Visit("https://github.com/settings/profile")
	return c
}

func hunterGetData(c1 *colly.Collector, keywords []string, page int) ([]string, []string) {
	var link []string
	var code []string
	c2 := c1.Clone()
	c2.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:66.0) Gecko/20100101 Firefox/66.0"
	c2.Async = true
	c2.Limit(&colly.LimitRule{
		DomainGlob:  "*github.*",
		Parallelism: 10,
	})
	c2.OnHTML("div.code-list-item.col-12.py-4.code-list-item-public ", func(e2 *colly.HTMLElement) {
		s, _ := e2.DOM.Html()
		if strings.Contains(s, `<div class="file-box blob-wrapper">`) {
			tLink, _ := e2.DOM.Find("div.flex-auto.min-width-0.col-10 > a:nth-child(2)").Attr("href")
			tLink = "https://github.com" + tLink
			tCode, _ := e2.DOM.Find("div.file-box.blob-wrapper").Html()
			tCode = strings.Replace(tCode, "<span class='text-bold'>", `<span style="color:red">`, -1)
			link = append(link, tLink)
			code = append(code, tCode)

		}
	})
	for _, keyword := range keywords {
		progressBar := fmt.Sprintf(`{{ blue "正在收集: %s"}} {{counters . | blue}} {{bar . | blue}} {{percent . | yellow}} {{speed . | yellow}}`, keyword)
		bar := pb.ProgressBarTemplate(progressBar).Start(page)
		for i := 1; i <= page; i++ {
			if err := c2.Visit("https://github.com/search?o=desc&p=" + strconv.Itoa(i) + "&q=" + keyword + "&s=indexed&type=Code"); err != nil {
				fmt.Printf("\x1b[0;31m[-]收集 %s 的第 %d 页失败！具体错误: %s\x1b[0m\n", keyword, i, err.Error())
				time.Sleep(1 * time.Second)
			} else {
				//fmt.Printf("\x1b[0;35m[*]正在收集 %s 的第 %d 页信息...\x1b[0m\n", keyword, i)
				bar.Increment()
				time.Sleep(1 * time.Second)
			}
		}
		bar.Finish()
		c2.Wait()
	}
	fmt.Println("\x1b[0;36m[√]所有信息收集完成！\x1b[0m")
	return link, code
}

func sendMail(receiver []string, sender, message, host, user, password string) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", sender, "Github信息泄露监控")
	m.SetHeaders(map[string][]string{
		"To":      receiver,
		"Subject": {"Github信息泄露通知"},
	})
	m.SetBody("text/html", message)

	d := gomail.NewDialer(host, 25, user, password)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("发送邮件失败！", err)
		panic(err)
	}
	fmt.Println("邮件发送成功！")
}

func insertDB(url, code string) {
	db, err := sql.Open("sqlite3", "hunter.db")
	checkError(err)
	_, erro := db.Exec("CREATE TABLE IF NOT EXISTS Baseline (url varchar(1000) PRIMARY KEY, code varchar(10000) UNIQUE)")
	checkError(erro)
	stmt, err := db.Prepare("INSERT OR REPLACE INTO Baseline (url, code) values (?,?)")
	checkError(err)
	_, er := stmt.Exec(url, code)
	checkError(er)
}

func compareDBURL(url string) (string, string, error) {
	var searchedURL string
	var searchedCode string
	db, err := sql.Open("sqlite3", "hunter.db")
	if err != nil {
		fmt.Println("数据库连接失败！")
		checkError(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT url,code from Baseline where url = ?", url)
	err = row.Scan(&searchedURL, &searchedCode)
	return searchedURL, searchedCode, err
}

func initInfomation() ([]string, []string) {
	var (
		urls  []string
		codes []string
	)
	db, err := sql.Open("sqlite3", "hunter.db")
	if err != nil {
		fmt.Println("数据库连接失败! ")
		checkError(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT url,code from Baseline")
	checkError(err)
	defer rows.Close()
	for rows.Next() {
		var (
			url  string
			code string
		)
		if err = rows.Scan(&url, &code); err != nil {
			log.Fatal(err)
		}
		urls = append(urls, url)
		codes = append(codes, code)
	}
	return urls, codes
}

func checkError(err error) {
	if err != nil {
		fmt.Println("出现错误，信息如下：")
		panic(err)
	}
}

func banner() {
	frontColor := "\x1b[0;36m"
	backColor := "\x1b[0m"
	banner := fmt.Sprintf(`%s
     #####                                  #     #                                   
    #     # # ##### #    # #    # #####     #     # #    # #    # ##### ###### #####  
    #       #   #   #    # #    # #    #    #     # #    # ##   #   #   #      #    # 
    #  #### #   #   ###### #    # #####     ####### #    # # #  #   #   #####  #    # 
    #     # #   #   #    # #    # #    #    #     # #    # #  # #   #   #      #####  
    #     # #   #   #    # #    # #    #    #     # #    # #   ##   #   #      #   #  
     #####  #   #   #    #  ####  #####     #     #  ####  #    #   #   ###### #    #    V2.1 Created by Allen 
	%s`, frontColor, backColor)
	fmt.Println(banner)
}
