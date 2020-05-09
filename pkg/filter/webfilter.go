package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bsir2020/basework/configs"
	"github.com/bsir2020/basework/pkg/auth"
	"github.com/bsir2020/basework/pkg/rsa"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Filter struct {
}

type loginModul struct {
	Uid int64 `json:"uid"`
}

func (f *Filter) buildResponse(code int, status bool, errmsg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":      code,
		"success":   status,
		"err_msg":   errmsg,
		"data":      data,
		"timestamp": time.Now().String(),
	})

	c.Abort()
}

//请求head,必须包含auth,exp项
func (f *Filter) Checkauth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("===========", c.FullPath())
		if _, ok := configs.WhiteList[c.FullPath()]; ok {
			//放行
			c.Next()

			return
		}

		jwt := auth.New()
		a := c.Request.Header.Get("auth")
		if a == "" {
			f.buildResponse(1001, false, "token为空", nil, c)
			return
		}
		//e := c.Request.Header.Get("exp")

		//isOK = true

		var err error
		//解密
		if configs.EnvConfig.RunMode != 1 {
			a, err = rsa.RsaDecrypt(a)
			if err != nil {
				f.buildResponse(1002, false, err.Error(), nil, c)
				return
			}
		}

		/*
			expData, err := rsa.RsaDecrypt(e)
			if err != nil {
				f.respondWithError(10002, err.Error(), c)
				return
			}

			//超时
			t, _ := strconv.ParseInt(string(expData), 10, 64)
			if time.Now().Unix() > t {
				f.respondWithError(10003, err.Error(), c)
				return
			}
		*/

		//token
		if !jwt.TokenIsInvalid(a) {
			f.buildResponse(1004, false, "token is invalid", nil, c)
			return
		}

		loginmodul := loginModul{}
		data, err := c.GetRawData()
		if err != json.Unmarshal(data, &loginmodul) {
			f.buildResponse(500, false, err.Error(), nil, c)
			return
		}

		fmt.Println(loginmodul)

		if m, err := jwt.ParseToken(a); err != nil {
			f.buildResponse(500, false, err.Error(), nil, c)
			return
		} else {
			uid := m["uid"].(string)
			Uid, _ := strconv.ParseInt(uid, 10, 64)
			fmt.Println("param--> ", loginmodul)
			if Uid != loginmodul.Uid {
				f.buildResponse(1006, false, "user is invalid", nil, c)
				return
			}
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		//放行
		c.Next()
	}
}
