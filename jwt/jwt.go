package jwt

import (
	"crypto/rsa"
	"errors"
	"flag"
	"io/ioutil"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qiaogw/pkg/genrsa"
)

//EasyToken token结构
type EasyToken struct {
	Appid  string `json:"appid"`
	Userid string `json:"userid"`
	jwt.StandardClaims
	Name    string `json:"username"`
	IsAdmin bool   `json:"isadmin"`
	// User      models.User
}

// 定义秘钥地址
const (
// privKeyPath = "key/private.pem" // 私钥
// pubKeyPath  = "key/public.pem"  //公钥

)

//定义jwt加解密秘钥
var (
	verifyKey    *rsa.PublicKey
	mySigningKey *rsa.PrivateKey
	// jwtPrefix    = "Bearer"
)

//genkey 检查秘钥是否存在，若不存在生成秘钥
func genkey(jwtPublicPath, jwtPrivatePath string) (err error) {
	var bits int
	flag.IntVar(&bits, "btghhhhh", 2048, "密钥长度，默认为1024位")
	err = genrsa.GenKey(bits, jwtPublicPath, jwtPrivatePath)
	return
}

//初始化jwt秘钥
func InitJwt(jwtPublicPath, jwtPrivatePath string) (err error) {
	verifyBytes, err := ioutil.ReadFile(jwtPublicPath)
	if err != nil {
		err = genkey(jwtPublicPath, jwtPrivatePath)
		if err != nil {
			return err
		}
		verifyBytes, err = ioutil.ReadFile(jwtPublicPath)
		if err != nil {
			return err
		}
	}
	//公钥解密
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	signBytes, err := ioutil.ReadFile(jwtPrivatePath)

	if err != nil {
		return err
	}
	//私钥加密
	mySigningKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	return
}

//GetToken 获取token
func (e EasyToken) GetToken(tokenExpire int64) (string, error) {
	tokenlife := time.Now().Unix() + tokenExpire
	// 创建 Claims
	claims := EasyToken{
		e.Appid,
		e.Userid,
		jwt.StandardClaims{
			ExpiresAt: tokenlife,
			Issuer:    e.Appid,
			IssuedAt:  time.Now().Unix(),
		},
		e.Name,
		e.IsAdmin,
		// e.User,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	//私钥签发
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

//RefreshGetToken 刷新token
func (e EasyToken) refreshGetToken(tokenExpire int64) (string, error) {
	if (time.Now().Unix() - e.IssuedAt) > (tokenExpire) {
		return "", errors.New("token创建已超过系统设置，不能刷新")
	}
	tokenlife := time.Now().Unix() + tokenExpire
	// Create the Claims
	claims := EasyToken{
		e.Appid,
		e.Userid,
		jwt.StandardClaims{
			ExpiresAt: tokenlife,
			Issuer:    e.Appid,
			IssuedAt:  e.IssuedAt,
		},
		e.Name,
		e.IsAdmin,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	//私钥签发
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

//ValidateToken 验证token
func ValidateToken(tokenString string) bool {
	if tokenString == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if token == nil {
		return false
	}

	if token.Valid {
		return true
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// 令牌已过期或尚未激活
			return false
		} else {
			//无法处理此令牌
			return false
		}
	} else {
		//无法处理此令牌
		return false
	}
}

//Parse 验证并获取token对象内容，若令牌过去则刷新
func Parse(tokenString string, tokenExpire int64) (*EasyToken, string, error) {
	if tokenString == "" {
		return nil, "", errors.New("Token为空！")
	}
	token, err := jwt.ParseWithClaims(tokenString, &EasyToken{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if token == nil {
		return nil, "", errors.New("token不工作")
	}
	if token.Valid {
		if claims, ok := token.Claims.(*EasyToken); ok {
			//令牌合法，解析成功，返回解析内容和现有令牌
			return claims, tokenString, nil
		}
		return nil, "", errors.New("token无法解析")

	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, "", errors.New("这不是一个合法的令牌！")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// 令牌已过期或尚未激活
			if claims, ok := token.Claims.(*EasyToken); ok {
				//开始刷新令牌
				ts, err := claims.refreshGetToken(tokenExpire)
				if err != nil {
					return nil, "", errors.New("token已超过刷新周期，重新登录！")
				}
				return claims, ts, nil
			}
			return nil, "", errors.New("token已过期")

		} else {
			//无法处理此令牌
			return nil, "", err
		}
	} else {
		//无法处理此令牌
		return nil, "", err
	}
}
