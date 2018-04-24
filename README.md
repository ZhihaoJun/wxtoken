# wxtoken

weixin access token and jsapi signature service

* help you to cache access_token and provide HTTP endpoint
* help you to cache jsapi_ticket
* help you to sign jssdk config



## Usage

pull down the image of `zhihaojun/wxtoken` , and run it

image environment variable config

* `WXTOKEN_APPID` app id
* `WXTOKEN_APPSECRET` app secret
* `WXTOKEN_ADDR` format of `<ip>:<port>` to make app listen on other port, default is `3001`



**The container don't have any authentication mechanism. DONT expose any api of the container on the public internet!** 



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
	"timestamp": 151234212
  }
}
```



`GET /ping`

response nothing in body with http code of 200



## Dependencies

* echo 2.0
* golang

