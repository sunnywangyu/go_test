package main

import (
	"fmt"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//巡检白名单
var whiteList = map[string]int{
	"test": 1,
	"old":1,
	"bak":1,
}
func main() {
	c := cron.New()
	//定时任务
	spec := "0 */1 * * * ?" //cron表达式，每分钟执行一次
	c.AddFunc(spec, func() {
		checkFile()
		fmt.Println("执行巡检任务===",time.Now().Format( "200601021504"))
	})

	c.Start()
	select {}  //阻塞主线程停止
}

func checkFile() {
	//循环目录下的文件夹
	feishuUrl := "https://ksc.webhook.ksyun.com/open-apis/bot/v2/hook/a2d15090-d825-4e60-bd7d-1197b7bc7a86"
	path := "/data/nfsen/profiles-data/live/"
	dirInfos,err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("读取子文件夹失败,path=",path,",err=",err)
	}
	var errorScore string
	duration,_ := time.ParseDuration("-5m")
	checkFileName :=  time.Now().Add(duration).Format( "200601021504")
	checkFilePath := time.Now().Add(duration).Format( "2006/01/02")
	for _,fi := range dirInfos{
		//filename := strings.Split(fi.Name(),"_")
		//endTag := filename[1]
		//if _,ok := whiteList[endTag];ok {//如果文件后缀在白名单中停止检查报警
		//	continue

		//}
		if _,ok := whiteList[fi.Name()];ok {//如果文件夹在白名单中停止检查报警
			continue
		}
		checkFile := path + fi.Name() + "/" + checkFilePath + "/nfcapd." + checkFileName
		_, err := os.Stat(checkFile)
		if err != nil {
			fmt.Println("获取文件失败,file=",checkFile,",err=",err)
			errorScore += `{
				"tag":"div",
					"text":{
					"content": "` + fi.Name() + `",
						"tag":"lark_md"
				}
			},`
		}
	}
	if len(errorScore) > 0 {
		errorScore = errorScore[:len(errorScore)-1]
		//发送到飞书报警群里
		stringData := `
		{
			"msg_type": "interactive",
			"card": {
				"config": {
					"wide_screen_mode":true
				},
                "header" : {
					"template" : "red",
					"title": {
						"content":"NFSEN服务巡检",
						"tag":"plain_text"
					}
                },
				"elements":[
					{
						"tag":"div",
                        "text":{
							"content": "巡检时间: ` + checkFileName + `，异常设备如下：", 
							"tag":"lark_md"
						}
                    },	
					{
 						"tag":"hr"
					},
					` + errorScore + `
				]
			}
		}
		`
		resp,err := http.Post(feishuUrl,"application/json", strings.NewReader(stringData))
		if err != nil {
			fmt.Println("调用飞书接口异常,error = ",err)
		} else {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
		}
	} else {
		//获取当前是否为整点
		hours, minutes, _ := time.Now().Clock()
		if strconv.Itoa(hours) == "12" && strconv.Itoa(minutes) == "30" {
			//每天12点发送一个验证执行中，成功的日志
			stringData := `
		{
			"msg_type": "interactive",
			"card": {
				"config": {
					"wide_screen_mode":true
				},
                "header" : {
					"template" : "green",
					"title":{
						"content":"NFSEN服务巡检",
						"tag":"plain_text"
					}
                },
				"elements":[ 
					{
						"tag":"div",
                        "text":{
							"content": "巡检时间: ` + checkFileName + `，无异常设备",
							"tag":"lark_md"
						}
                    }	
				]
			}
		}
		`
			resp, err := http.Post(feishuUrl,"application/json", strings.NewReader(stringData))
			if err != nil {
				fmt.Println("调用飞书接口异常,error = ",err)
			} else {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(body))
			}

		}
	}

	return
}
