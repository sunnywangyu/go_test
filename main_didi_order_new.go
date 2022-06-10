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
	getDiDiUserCarConfig()
	/*c := cron.New()
	spec := "0 0 9 * * ?" //cron表达式，每天9点执行一次
	c.AddFunc(spec, func() {
		clientId = "1cc34da2484e52b8921d398298909abb"
		clientSecret = "bab21ae4b9820bc62dc0b9c69e80813e"
		grantType = "client_credentials"
		phone = "13552847730"
		signkey = "92794F49b7efEF94737c"
	    companyId = "1125909343514202"
		getDiDiCarInfo()
		fmt.Println("调用滴滴接口入库完成===",time.Now().Format( "200601021504"))
	})
	c.Start()
	//阻塞主线程停止
	select {}*/
}

//获取滴滴用车规则
func getDiDiUserCarConfig(){
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
		json.Unmarshal(body, &RespAccessToken)
		if RespAccessToken.Errno > 0 {
			fmt.Println("调用滴滴获取access_token接口异常,error = ", RespAccessToken.Errmsg)
			return
		}
		//解析access_token
		accessToken := RespAccessToken.Token
		//拼装用车规则参数
		orderUrlParam := "access_token=" + accessToken + "&client_id=" + clientId + "&company_id=" + companyId +
			"&sign_key=" + signkey +  "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10)
		signOrder :=  fmt.Sprintf("%x", md5.Sum([]byte(orderUrlParam)))
		//调用获取订单接口
		orderUrl := "https://api.es.xiaojukeji.com/river/Regulation/get?"+ orderUrlParam + "&sign="+ signOrder
		respOrder,err2 := http.Get(orderUrl)
		if err2 != nil {
			fmt.Println("调用滴滴获取用车规则接口异常,error = ",err2)
		} else {
			defer respOrder.Body.Close()
			//解析订单json
			respOrderCarConfigBody, _ := ioutil.ReadAll(respOrder.Body)
			respOrderCarConfigBodyString := string(respOrderCarConfigBody)
			var DiDiOrderCarConfig map[string]interface{}
			json.Unmarshal([]byte(respOrderCarConfigBodyString), &DiDiOrderCarConfig)
			if DiDiOrderCarConfig["errno"].(float64) > 0 {
				fmt.Println("调用滴滴获取用车规则接口异常,error = ", DiDiOrderCarConfig["errmsg"])
			}
			//total := DiDiOrderCarConfig["data"].(map[string]interface{})["total"].(float64)
			fmt.Println(DiDiOrderCarConfig)
		}
	}
}

func getDiDiUserInfo(){
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
		json.Unmarshal(body, &RespAccessToken)
		if RespAccessToken.Errno > 0 {
			fmt.Println("调用滴滴获取access_token接口异常,error = ", RespAccessToken.Errmsg)
			return
		}
		//解析access_token
		accessToken := RespAccessToken.Token
		//拼装订单参数
		orderUrlParam := "access_token=" + accessToken + "&client_id=" + clientId + "&company_id=" + companyId +
			"&length=" +  strconv.Itoa(10) + "&offset=" + strconv.Itoa(0) +
			"&sign_key=" + signkey +  "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10)
		signOrder :=  fmt.Sprintf("%x", md5.Sum([]byte(orderUrlParam)))
		//调用获取订单接口
		orderUrl := "https://api.es.xiaojukeji.com/river/Member/get?"+ orderUrlParam + "&sign="+ signOrder
		respOrder,err2 := http.Get(orderUrl)
		if err2 != nil {
			fmt.Println("调用滴滴获取订单接口异常,error = ",err2)
		} else {
			defer respOrder.Body.Close()
			//解析订单json
			respOrderBody, _ := ioutil.ReadAll(respOrder.Body)
			respOrderBodyString := string(respOrderBody)
			var DiDiOrder map[string]interface{}
			json.Unmarshal([]byte(respOrderBodyString), &DiDiOrder)
			if DiDiOrder["errno"].(float64) > 0 {
				fmt.Println("调用滴滴获取订单接口异常,error = ", DiDiOrder["errmsg"])
			}
			//total := DiDiOrder["data"].(map[string]interface{})["total"].(float64)
			fmt.Println(DiDiOrder)
		}
	}
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
			for i := 1; i < totalPageInt; i++ {
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
		reponsePush,err3 := http.Post(pushOrderUrl,"application/json", strings.NewReader(postOrderData))
		if err3 != nil {
			fmt.Println("推送滴滴订单数据到BI库,error = ",err3)
		} else {
			defer respOrder.Body.Close()
			//解析推送结果
			reponsePushBody, _ := ioutil.ReadAll(reponsePush.Body)
			var pushResult map[string]interface{}
			json.Unmarshal([]byte(string(reponsePushBody)),&pushResult)
			fmt.Println(pushResult)
			if pushResult["status"].(bool) == true && pushResult["total"].(float64) > 0 {
				for _,v := range pushResult["msg"].([]interface{}) {
					v = v.(map[string]interface{})["member_id"]
					 //调用滴滴接口获取用户信息详情
					//拼装用户详情参数
					userUrlParam := "access_token=" + accessToken + "&client_id=" + clientId + "&company_id="+ companyId +
						"&member_id=" + v.(string) +"&sign_key=" + signkey +  "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10)
					signUser :=  fmt.Sprintf("%x", md5.Sum([]byte(userUrlParam)))
					//调用获取订单接口
					userUrl := "https://api.es.xiaojukeji.com/river/Member/detail?"+ userUrlParam + "&sign="+ signUser
					userInfo,err4 := http.Get(userUrl)
					if err4 != nil {
						fmt.Println("调用滴滴获取用户详情接口异常,member_id = ",v,",error=",err4)
						continue
					} else {
						defer userInfo.Body.Close()
						//解析订单json
						userInfoBody, _ := ioutil.ReadAll(userInfo.Body)
						var DiDiUserInfo map[string]interface{}
						json.Unmarshal([]byte(string(userInfoBody)),&DiDiUserInfo)
						if DiDiUserInfo["errno"].(float64) > 0 {
							fmt.Println("调用滴滴获取用户详情接口异常,member_id =",v,",error = ",DiDiUserInfo["errmsg"],",data=",DiDiUserInfo["data"])
							continue
						}
						//将订单records数据转成json
						pushUserInfo,_ := json.Marshal(DiDiUserInfo["data"].(map[string]interface{}))
						//推送用户数据到Cost库参数
						pushUserInfoData := `
							{
								"access_token": "f8e947c7e7186e8d4a4cf9a0643b9a27",
								"table": "pull_didi_user_info",
								"tag": "` + tag + `",
								"data": ` + string(pushUserInfo) + `
							}
 						`
						//推送用户数据到Cost接口入库
						_,err5 := http.Post(pushOrderUrl,"application/json", strings.NewReader(pushUserInfoData))
						if err5 != nil {
							fmt.Println("推送滴滴用户数据到BI库,error = ",err5)
						}
					}
				}
			}
		}
		return total
	}

}
