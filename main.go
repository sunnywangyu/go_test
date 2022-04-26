package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AccessTokenStruct struct {
	Token string `json:"access_token"`
	Errno int `json:"errno"`
	Errmsg string `json:"errmsg"`
}

var (
	clientId string
	clientSecret string
	grantType string
	phone string
	signkey string
	companyId string
)
func main() {
	//c := cron.New()
	//定时任务
	//spec := "0 */5 * * * ?" //cron表达式，每分钟执行一次
	//c.AddFunc(spec, func() {
		clientId = "1cc34da2484e52b8921d398298909abb"
		clientSecret = "bab21ae4b9820bc62dc0b9c69e80813e"
		grantType = "client_credentials"
		phone = "13552847730"
		signkey = "92794F49b7efEF94737c"
	    companyId = "1125909343514202"
		getDiDiCarInfo()
		fmt.Println("调用滴滴接口入库完成===",time.Now().Format( "200601021504"))
	//})

	//c.Start()
	//select {}  //阻塞主线程停止
}

func getDiDiCarInfo() {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	tokenSignString := "client_id=" + clientId + "&client_secret=" + clientSecret + "&grant_type=" +
		grantType + "&phone=" + phone + "&sign_key=" + signkey + "&timestamp=" + timestamp
	//access_token参数
	stringData := `
		{
			"client_id": "` + clientId + `",
			"client_secret": "` + clientSecret + `",
			"grant_type": "` + grantType + `",
			"phone": "` + phone + `",
			"timestamp": "` + timestamp + `",
			"sign": "` + fmt.Sprintf("%x", md5.Sum([]byte(tokenSignString))) + `"
		}
		`
	//调用滴滴获取access_token接口
	accessTokenUrl := "https://api.es.xiaojukeji.com/river/Auth/authorize"
	resp,err := http.Post(accessTokenUrl,"application/json", strings.NewReader(stringData))
	if err != nil {
		fmt.Println("调用滴滴获取access_token接口异常,error = ",err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var RespAccessToken AccessTokenStruct
		json.Unmarshal(body,&RespAccessToken)
		if RespAccessToken.Errno > 0 {
			fmt.Println("调用滴滴获取access_token接口异常,error = ",RespAccessToken.Errmsg)
			return
		}
		//解析access_token
		accessToken := RespAccessToken.Token
		start_date_duration,_ := time.ParseDuration("-24h")
		start_date := time.Now().Add(start_date_duration).Format( "2006-01-02")
		//end_date := time.Now().Format( "2006-01-02")
		end_date := start_date
		tag := time.Now().Format( "200601021504") //每次调用算一批次
		pageSize := 100
		total := pushOrder(accessToken,end_date,start_date,tag,pageSize,0)
		totalPage := math.Ceil(total/float64(pageSize))
		totalPageInt, _ := strconv.Atoi(fmt.Sprintf("%1.0f",totalPage))
		if totalPageInt > 1 {
			for i := 0; i < totalPageInt; i++ {
				offset := i * pageSize
				//调用接口循环入库
				pushOrder(accessToken,end_date,start_date,tag,pageSize,offset)
			}
		}
	}
	return
}

func pushOrder(accessToken,end_date,start_date,tag string,pageSize,offset int )  float64 {
	//拼装订单参数
	orderUrlParam := "access_token=" + accessToken + "&client_id=" + clientId + "&company_id=" + companyId +
		"&end_date=" + end_date + "&length=" +  strconv.Itoa(pageSize) + "&offset=" + strconv.Itoa(offset) +
		"&sign_key=" + signkey + "&start_date=" + start_date +  "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10)
	signOrder :=  fmt.Sprintf("%x", md5.Sum([]byte(orderUrlParam)))
	//调用获取订单接口
	orderUrl := "http://api.es.xiaojukeji.com/river/Order/get?"+ orderUrlParam + "&sign="+ signOrder
	respOrder,err2 := http.Get(orderUrl)
	if err2 != nil {
		fmt.Println("调用滴滴获取订单接口异常,error = ",err2)
		return float64(0)
	} else {
		defer respOrder.Body.Close()
		//解析订单json
		respOrderBody, _ := ioutil.ReadAll(respOrder.Body)
		respOrderBodyString := string(respOrderBody)
		var DiDiOrder map[string]interface{}
		json.Unmarshal([]byte(respOrderBodyString),&DiDiOrder)
		if DiDiOrder["errno"].(float64) > 0 {
			fmt.Println("调用滴滴获取订单接口异常,error = ",DiDiOrder["errmsg"])
			return float64(0)
		}
		total := DiDiOrder["data"].(map[string]interface{})["total"].(float64)

		if total <= 0 {
			fmt.Println("调用订单接口,没有获取到订单数据,total = ",total)
			return float64(0)
		}
		fmt.Println("调用订单接口,获取到订单数据,total = ",total,",pageSize=",pageSize,",offset=",offset)
		//将订单records数据转成json
		jsonRecords,_ := json.Marshal(DiDiOrder["data"].(map[string]interface{})["records"])
		//推送订单数据到Cost库参数
		postOrderData := `
		{
			"access_token": "f8e947c7e7186e8d4a4cf9a0643b9a27",
			"table": "pull_didi_car_order",
			"tag": "` + tag + `",
			"data": ` + string(jsonRecords) + `
		}
		`
		//fmt.Println(postOrderData)
		pushOrderUrl := "http://120.92.106.81:1502/interface/bi-accept-data"
		//推送订单数据到Cost接口入库
		_,err3 := http.Post(pushOrderUrl,"application/json", strings.NewReader(postOrderData))
		if err3 != nil {
			fmt.Println("推送滴滴订单数据到BI库,error = ",err3)
		}
		return total
	}

}
