package filter

import (
	"github.com/bsir2020/basework/configs"
	"github.com/bsir2020/basework/pkg/auth"
	"github.com/bsir2020/basework/pkg/rsa"
	"github.com/gin-gonic/gin"
	"net/http"
	//"strconv"
	"time"
)

type Filter struct {
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

		//解密
		authData, err := rsa.RsaDecrypt(a)
		if err != nil {
			f.buildResponse(1002, false, err.Error(), nil, c)
			return
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
		if !jwt.TokenIsInvalid(string(authData)) {
			f.buildResponse(1004, false, err.Error(), nil, c)
			return
		}

		u, _ := c.GetPostForm("uid")

		if m, err := jwt.ParseToken(a); err != nil {
			f.buildResponse(1005, false, err.Error(), nil, c)
			return
		} else {
			if u != m["uid"] {
				f.buildResponse(1006, false, err.Error(), nil, c)
				return
			}
		}

		//放行
		c.Next()
	}
}
