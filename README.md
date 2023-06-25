<!--
title: .'HTTP GET and POST'
description: 'Boilerplate code for Golang with GET and POST example'
framework: v1
platform: AWS
language: Go
priority: 10
authorLink: 'https://github.com/pramonow'
authorName: 'Pramono Winata'
authorAvatar: 'https://avatars0.githubusercontent.com/u/28787057?v=4&s=140'
-->

# Serverless-golang http Get and Post Example
Serverless boilerplate code for golang with GET and POST example

This example is using AWS Request and Response Proxy Model, provided by AWS itself.
If you want to test any changes don't forget to run `make` inside the service directory.
 


### POST endpoint with name in the body /post 
- input đầu vào có 2 field: value1, value2
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "value1": {{number}},
        "value2": {{number}}
    }
}
```
- output: sẽ trả về giá trị của value1+value2
```
{
    "requestId": {{requestId}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "sum": {{value1+value2}}
    }
}
```
### POST endpoint with name in the body /postapi2 
- input đầu vào có 2 field: plaintText, secretKey
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "plaintText": {{string}},
        "secretKey": {{string}}
    }
}
```
- output: sẽ trả về 1 field: signature sử dụng thuật toán sha256 hoặc hmacsha256
```
{
    "requestId": {{requestId}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "signature": {{string}}
    }
}
```
### POST endpoint with name in the body /postapi3 
- là dùng base64, input có 2 filed: needEncode, needDecode
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "needEncode": {{string}},
        "needDecode": {{string}}
    }
}
```
- output: sẽ trả về 2 field: outEncode là output của base64 field needEncode, outDecode là output của field needDecode
```
{
    "requestId": {{requestId}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "outEncode": {{string}},
        "outDecode": {{string}}
    }
}
```