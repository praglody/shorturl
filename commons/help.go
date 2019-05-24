package commons

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const Success = 200      //正常
const Failed = 500       //失败
const ParamsError = 4001 //参数错误
const NotFound = 4004    //记录不存在

//打印元素
func DD(obj ...interface{}) {
	fmt.Printf("%v\n", obj)
}

//md5
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
