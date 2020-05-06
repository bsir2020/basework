package auth

import (
	"errors"
	"fmt"
	cfg "github.com/bsir2020/basework/configs"
	"github.com/bsir2020/basework/pkg/log"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

type JWT struct {
	signingKey string
	subject    string //主题
}

var authLog = log.New()

type Token struct {
	Token string `json:"token"`
}

func New() (jwt *JWT) {
	jwt = &JWT{
		signingKey: cfg.EnvConfig.Authkey.Key,
		subject:    cfg.EnvConfig.Authkey.Subject,
	}

	return
}

func (j *JWT) CreateToken(userid int, exptime int64) (res Token) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	//claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix() //过期时间
	claims["exp"] = exptime           //过期时间
	claims["iat"] = time.Now().Unix() //签发时间
	claims["sub"] = j.subject         //主题
	claims["uid"] = userid
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		fmt.Print("Error while signing the token")
		authLog.Fatal("CreateToken", zap.String("Error while signing the token", err.Error()))
	}

	res = Token{tokenString}
	return
}

func (j *JWT) ParseToken(tokenString string) (jwt.MapClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parse token err %v", token.Header["alg"])
		}
		return []byte(j.signingKey), nil
	})
	if err != nil {
		authLog.Error("ParseToken", zap.String("parse token error", err.Error()))
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

//if token is invalid, method will return true
func (j *JWT) TokenIsInvalid(tokenString string) bool {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		authLog.Fatal("TokenIsInvalid", zap.String("valid token error", err.Error()))
	} else {
		//校验下token是否过期
		if res := claims.VerifyExpiresAt(time.Now().Unix(), true); res == false {
			return true
		}

		if res := claims.VerifyIssuedAt(time.Now().Unix(), true); res == false {
			return true
		}

		if res := claims["sub"].(string); res == j.subject {
			return false
		}

		if res := claims["uid"].(int); res == 0 {
			return true
		}
	}

	return true
}
