package jwt

import (
	"crypto/rsa"
	"errors"
	"flag"
	"io/ioutil"

	"github.com/qiaogw/pkg/config"

	//"github.com/qiaogw/pkg/conf"
	//"github.com/qiaogw/pkg/config"
	"time"

	"github.com/qiaogw/pkg/genrsa"
	"github.com/qiaogw/pkg/logs"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
)

//EasyToken token结构
type EasyToken struct {
	Appid  string `json:"appid"`
	Userid int    `json:"userid"`
	Orgid  int    `json:"orgid"`
	jwt.StandardClaims
	Name      string   `json:"username"`
	OrgName   string   `json:"orgname"`
	RoleIds   []int    `json:"roleids"`
	RoleNames []string `json:"rolename"`
	IsAdmin   bool     `json:"isadmin"`
	//User      models.User
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
func genkey() (err error) {
	var bits int
	flag.IntVar(&bits, "btghhhhh", 2048, "密钥长度，默认为1024位")
	err = genrsa.GenKey(bits)
	return
}

//初始化jwt秘钥
func InitJwt() {
	verifyBytes, err := ioutil.ReadFile(config.Config.JwtPublicPath)
	if err != nil {
		beego.Error(err, config.Config.JwtPublicPath)
		err = genkey()
		if err != nil {
			logs.Fatal(err)
		}
		verifyBytes, err = ioutil.ReadFile(config.Config.JwtPublicPath)
		if err != nil {
			beego.Error(err)
		}
	}
	//公钥解密
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		beego.Error(err)
	}

	signBytes, err := ioutil.ReadFile(config.Config.JwtPrivatePath)

	if err != nil {
		logs.Fatal(err)
	}
	//私钥加密
	mySigningKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		logs.Fatal(err)
	}
}

//GetToken 获取token
func (e EasyToken) GetToken() (string, error) {
	//maxlife, err := beego.AppConfig.Int64("sessiongcmaxlifetime")
	maxlife := config.Config.TokenExpire
	//if err != nil {
	//	logs.Error("获取配置错误！ err is ", err.Error())
	//	maxlife = 3600
	//}
	tokenlife := time.Now().Unix() + maxlife
	// 创建 Claims
	claims := EasyToken{
		e.Appid,
		e.Userid,
		e.Orgid,
		jwt.StandardClaims{
			ExpiresAt: tokenlife,
			Issuer:    e.Appid,
			IssuedAt:  time.Now().Unix(),
		},
		e.Name,
		e.OrgName,
		e.RoleIds,
		e.RoleNames,
		e.IsAdmin,
		//e.User,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	//私钥签发
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		logs.Error(err)
		return tokenString, err
	}
	logs.Info(e.Appid, "为用户 ", e.Userid, "创建了token！")
	return tokenString, nil
}

//RefreshGetToken 刷新token
func (e EasyToken) RefreshGetToken() (string, error) {
	refreshlife := config.Config.TokenExpire
	//判断token创建时间至今是否已超过强制登录时间
	if (time.Now().Unix() - e.IssuedAt) > (refreshlife) {
		return "", errors.New("token创建已超过系统设置，不能刷新")
	}
	maxlife := config.Config.TokenExpire
	tokenlife := time.Now().Unix() + maxlife
	// Create the Claims
	claims := EasyToken{
		e.Appid,
		e.Userid,
		e.Orgid,
		jwt.StandardClaims{
			ExpiresAt: tokenlife,
			Issuer:    e.Appid,
			IssuedAt:  e.IssuedAt,
		},
		e.Name,
		e.OrgName,
		e.RoleIds,
		e.RoleNames,
		e.IsAdmin,
		//e.User,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	//私钥签发
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		logs.Fatal(err)
		return tokenString, err
	}
	logs.Info(e.Appid, "为用户 ", e.Userid, "刷新了token！")
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
		logs.Error(err)
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
func Parse(tokenString string) (*EasyToken, string, error) {
	if tokenString == "" {
		return nil, "", errors.New("Token为空！")
	}
	token, err := jwt.ParseWithClaims(tokenString, &EasyToken{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if token == nil {
		logs.Error(err)
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
				ts, err := claims.RefreshGetToken()
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
