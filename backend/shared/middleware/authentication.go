package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/authentication"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/runtimebag"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	PubKeyPath = "./public_key.pem"
)

var verifyKey *rsa.PublicKey

// Authenticate middleware to authenticate user
func Authenticate() gin.HandlerFunc {
	addJwtMiddleware()
	return func(ctx *gin.Context) {
		bToken := ctx.GetHeader("Authorization")
		if bToken == "" {
			bToken = ctx.Query("Authorization")
			if bToken == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("missing bearer token"))
				return
			}
		}
		if strings.Contains(bToken, "Bearer") {
			bToken = strings.Split(bToken, " ")[1]
		}
		token, err := jwt.ParseWithClaims(bToken, &authentication.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("invalid token. %s", err.Error()))
			return
		}
		ctx.Set("user", token)
		ctx.Next()
	}
}

type User interface {
	GetClaims() *authentication.CustomClaims
	IsAgent() bool
}

type AppUser struct {
	Claims *authentication.CustomClaims `json:"custom_claims"`
}

func (u *AppUser) IsAgent() bool {
	if u.Claims.UserID != "" {
		return true
	}
	return false
}

func (u *AppUser) GetClaims() *authentication.CustomClaims {
	return u.Claims
}

func BuildUser(tokenString string) (User, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &authentication.CustomClaims{})
	if err != nil {
		return nil, err
	}
	return &AppUser{
		Claims: token.Claims.(*authentication.CustomClaims),
	}, nil
}

func addJwtMiddleware() {
	keyLocation := runtimebag.GetEnvString(constants.PubKeySaveLocation, PubKeyPath)
	verifyBytes, err := os.ReadFile(keyLocation)
	if err != nil {
		log.Fatal(err.Error())
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
}

type RealmData struct {
	Realm          string `json:"realm"`
	PublicKey      string `json:"public_key"`
	TokenService   string `json:"token-service"`
	AccountService string `json:"account-service"`
}

func FetchKey() {
	url := runtimebag.GetEnvString(constants.IamUri, "")
	if url == "" {
		log.Fatalf("no %s defined in enviroment", constants.IamUri)
	}
	realm := runtimebag.GetEnvString(constants.IamRealm, "")
	if realm == "" {
		log.Fatalf("no %s defined in enviroment", constants.IamRealm)
	}
	url = fmt.Sprintf("%s/auth/realms/%s", url, realm)
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err.Error())
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	rd := &RealmData{}
	err = json.Unmarshal(b, rd)
	if err != nil {
		log.Fatal(err.Error())
	}
	var target *os.File

	keyLocation := runtimebag.GetEnvString(constants.PubKeySaveLocation, PubKeyPath)

	target, err = os.Create(keyLocation)
	if err != nil {
		return
	}
	if _, err = target.WriteString("-----BEGIN CERTIFICATE-----" + "\n"); err != nil {
		log.Fatal(err.Error())
	}
	if _, err = target.WriteString(rd.PublicKey + "\n"); err != nil {
		log.Fatal(err.Error())
	}
	if _, err = target.WriteString("-----END CERTIFICATE-----"); err != nil {
		log.Fatal(err.Error())
	}
	addJwtMiddleware()
	if err = target.Close(); err != nil {
		log.Panic(err)
	}
}
