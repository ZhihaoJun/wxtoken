# wxtoken

weixin access token and jsapi signature service

* help you to cache access_token and provide HTTP endpoint
* help you to cache jsapi_ticket
* help you to sign jssdk config



## usage

pull down the image of `zhihaojun/wxtoken` , and run

image environment variable config

* `WXTOKEN_APPID` app id
* `WXTOKEN_APPSECRET` app secret
* `WXTOKEN_ADDR` format of `<ip>:<port>` to make app listen on other port



**The container doesn't have any authentication mechanism. DONT expose any api of the container on the public internet!** 



## API
### access token

`GET /access_token `

response

```json
{
  "t": "<access_token>"
}
```


### jsapi ticket

`GET /jsapi_ticket`

response

```json
{
  "t": "<access_token>"
}
```

### jssdk config

`GET /jssdk_config`

query string parameters

* url: the page full url you are requesting
  * should be url entity escaped
  * the content after # will be stripped out automatically

response

```json
{
  "error": "ok",
  "msg": "get jssdk config success",
  "config": {
    "appId": "<appid>",
	"nonceStr": "<noncestr>",
	"signature": "<signature>",
	"timestamp": 151234212
  }
}
```



`GET /ping`

response nothing in body with http code of 200



## dependencies

* echo 2.0
* golang

