package filter

import (
	"github.com/bsir2020/basework/pkg/auth"
	"github.com/bsir2020/basework/pkg/rsa"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type Filter struct {
}

func (f *Filter) respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}

//请求head,必须包含auth,exp项
func (f *Filter) Checkauth() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt := auth.New()
		a := c.Request.Header.Get("auth")
		e := c.Request.Header.Get("exp")

		//isOK = true

		//解密
		authData, err := rsa.RsaDecrypt([]byte(a))
		if err != nil {
			f.respondWithError(10001, err.Error(), c)
			return
		}

		expData, err := rsa.RsaDecrypt([]byte(e))
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

		//token
		if !jwt.TokenIsInvalid(string(authData)) {
			f.respondWithError(10004, err.Error(), c)
			return
		}

		u, _ := c.GetPostForm("uid")

		if m, err := jwt.ParseToken(a); err != nil {
			f.respondWithError(10005, err.Error(), c)
			return
		} else {
			if u != m["uid"] {
				f.respondWithError(10006, "uid no match", c)
				return
			}
		}

		//放行
		c.Next()
	}
}
