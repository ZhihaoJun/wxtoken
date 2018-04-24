# wxtoken

weixin access token and jsapi signature service

* help you to cache access_token and provide HTTP endpoint
* help you to cache jsapi_ticket
* help you to signature jssdk config



## usage

pull down the image of `zhihaojun/wxtoken` , and run

image environment config

* `WXTOKEN_APPID` app id
* `WXTOKEN_APPSECRET` app secret
* `WXTOKEN_ADDR` format of `<ip>:<port>` to make app listen on other port



**This container don't have any authentication mechanism, DONT expose any api of the container on the internet!** 



## API

`GET /access_token `

response

```json

```



`GET /jsapi_ticket`

response

```json

```



`GET /jssdk_config`

query string parameters

* url: the page full url you are requesting
  * should be url entity escaped
  * the content after # should be stripped out by user

response

```json
{
  "error": "ok",
  "msg": "get jssdk config success",
  "config": {
    "appId": "<appid>",
	"nonceStr": "<noncestr>",
	"signature": "<signature>",
	"timestamp": <timestamp:int>
  }
}
```



`GET /ping`

response nothing in body with http code of 200



## dependencies

* echo 2.0
* golang

