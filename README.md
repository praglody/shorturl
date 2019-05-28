# shorturl
A Short Url Service Written By Golang

## 特性

- 使用 mysql 存储原始数据
- 使用lru cache 缓存热点数据，提高查询速度
- 支持点击数统计，统计数据异步落盘

## 介绍

主要使用golang的gin框架，日志采用的beego的组件，数据库采用的gorm，为便于开发，目录结构被设计成类似PHP Web项目,目前主要实现了3个接口，包括网址生成、短网址的查询，1个批量生成接口，还有一个页面用于短网址的解析跳转功能！

### **1.短网址生成**
###### URL
```
/v1/query
```
###### Content-Type
```
application/x-www-form-urlencoded
```
###### HTTP请求方式
```
POST
```
###### 请求参数
|参数|必选|类型|说明|
|:-----  |:-------|:-----|-----                               |
|url    |true    |string|用来生成短网址的url                          |

###### 响应结果
```json
{
    "code": 200,
    "data": {
        "url": "http://127.0.0.1:8080/mCv"
    },
    "msg": "ok"
}
```

### **2.短网址查询**
###### URL
```
/v1/query
```
###### Content-Type
```
application/x-www-form-urlencoded
```
###### HTTP请求方式
```
POST
```
###### 请求参数
|参数|必选|类型|说明|
|:-----  |:-------|:-----|-----                               |
|url    |true    |string|生成的短网址                          |

###### 响应结果
```json
{
    "code": 200,
    "data": {
        "url": "http://www.google.com"
    },
    "msg": "ok"
}
```
### **3.批量生成短网址**
###### URL
```
/v1/multicreate
```
###### Content-Type
```
application/json
```
###### HTTP请求方式
```
POST
```
###### 请求参数
|参数|必选|类型|说明|
|:-----  |:-------|:-----|-----                               |
|url    |true    |array|多个需生成短网址的url,最多50个                          |

###### 响应结果
```json
{
    "code": 200,
    "data": {
        "urls": {
            "http://www.abc.com": "XrUnl",
            "http://www.baidu.com": "ArCnl"
        }
    },
    "msg": "ok"
}
```
