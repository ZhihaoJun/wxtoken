package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type weixinAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

var appid string
var tokenMutex sync.RWMutex
var token = ""
var tokenCH chan string
var jsapiTicketMutex sync.RWMutex
var jsapiTicket = ""
var accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
var jsapiTicketURL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getAccessTokenURL(appid, secret string) string {
	return fmt.Sprintf(accessTokenURL, appid, secret)
}

func WXConfigSign(jsapiTicket, nonceStr, timestamp, url string) (signature string) {
	if i := strings.IndexByte(url, '#'); i >= 0 {
		url = url[:i]
	}

	n := len("jsapi_ticket=") + len(jsapiTicket) +
		len("&noncestr=") + len(nonceStr) +
		len("&timestamp=") + len(timestamp) +
		len("&url=") + len(url)
	buf := make([]byte, 0, n)

	buf = append(buf, "jsapi_ticket="...)
	buf = append(buf, jsapiTicket...)
	buf = append(buf, "&noncestr="...)
	buf = append(buf, nonceStr...)
	buf = append(buf, "&timestamp="...)
	buf = append(buf, timestamp...)
	buf = append(buf, "&url="...)
	buf = append(buf, url...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}

// views
func pingView(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func accessTokenView(c echo.Context) error {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"t": token,
	})
}

func requestAccessToken(appid, secret string) (*weixinAccessTokenResponse, error) {
	resp, err := http.Get(getAccessTokenURL(appid, secret))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respBody weixinAccessTokenResponse
	err = json.Unmarshal(content, &respBody)
	if respBody.Errcode != 0 {
		return &respBody, errors.New(fmt.Sprintf("%d: %s", respBody.Errcode, respBody.Errmsg))
	}
	return &respBody, err
}

func cacheAccessToken(appid, secret string) (string, int) {
	log.Println("requesting new token")
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	resp, err := requestAccessToken(appid, secret)
	if err != nil {
		log.Println(err)
		return "", 1
	}
	token = resp.AccessToken
	log.Println(resp)
	return resp.AccessToken, resp.ExpiresIn
}

type Config struct {
	AppID     string
	AppSecret string
	Addr      string
}

func getConfig() *Config {
	c := &Config{
		AppID:     "",
		AppSecret: "",
		Addr:      ":3001",
	}

	if v, ok := os.LookupEnv("WXTOKEN_APPID"); ok {
		c.AppID = v
	}

	if v, ok := os.LookupEnv("WXTOKEN_APPSECRET"); ok {
		c.AppSecret = v
	}

	if v, ok := os.LookupEnv("WXTOKEN_ADDR"); ok {
		c.Addr = v
	}

	return c
}

type weixinJSApiTicketResponse struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
}

func jsapiTicketView(c echo.Context) error {
	jsapiTicketMutex.RLock()
	defer jsapiTicketMutex.RUnlock()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"t": jsapiTicket,
	})
}

func jssdkConfigView(c echo.Context) error {
	url := c.QueryParam("url")

	jsapiTicketMutex.RLock()
	defer jsapiTicketMutex.RUnlock()

	nonceStr := RandStringRunes(32)
	timestamp := time.Now().Unix()
	timestampStr := fmt.Sprintf("%d", timestamp)
	sign := WXConfigSign(jsapiTicket, nonceStr, timestampStr, url)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"error": "ok",
		"msg":   "get jssdk config success",
		"config": map[string]interface{}{
			"appId":     appid,
			"nonceStr":  nonceStr,
			"signature": sign,
			"timestamp": timestamp,
		},
	})
}

func getJSApiTicketURL(accessToken string) string {
	return fmt.Sprintf(jsapiTicketURL, accessToken)
}

func requestJSApiTicket(accessToken string) (*weixinJSApiTicketResponse, error) {
	resp, err := http.Get(getJSApiTicketURL(accessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respBody weixinJSApiTicketResponse
	err = json.Unmarshal(content, &respBody)
	if respBody.Errcode != 0 {
		return &respBody, errors.New(fmt.Sprintf("%d: %s", respBody.Errcode, respBody.Errmsg))
	}
	return &respBody, err
}

func cacheJSApiTicket(accessToken string) (string, int) {
	log.Println("requesting new jsapi ticket")
	jsapiTicketMutex.Lock()
	defer jsapiTicketMutex.Unlock()
	resp, err := requestJSApiTicket(accessToken)
	if err != nil {
		log.Println(err)
		return "", 1
	}
	jsapiTicket = resp.Ticket
	log.Println(resp)
	return resp.Ticket, resp.ExpiresIn
}

func logSkipper(c echo.Context) bool {
	return c.Request().URI() == "/ping"
}

func main() {
	c := getConfig()
	appid = c.AppID
	appsecret := c.AppSecret
	tokenCH = make(chan string)
	go func() {
		log.Println("start caching token")
		for {
			accessToken, t := cacheAccessToken(appid, appsecret)
			tokenCH <- accessToken
			time.Sleep(time.Second * time.Duration(t))
		}
	}()

	go func() {
		log.Println("start caching jsapi ticket")
		for {
			_, t := cacheJSApiTicket(<-tokenCH)
			time.Sleep(time.Second * time.Duration(t))
		}
	}()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: logSkipper,
	}))
	e.Use(middleware.Recover())
	e.GET("/access_token", accessTokenView)
	e.GET("/jsapi_ticket", jsapiTicketView)
	e.GET("/jssdk_config", jssdkConfigView)
	e.GET("/ping", pingView)
	e.Run(standard.New(c.Addr))
}
