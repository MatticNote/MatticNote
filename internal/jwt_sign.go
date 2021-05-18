package internal

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/google/uuid"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

var (
	jwtSignPublicKey  *rsa.PublicKey
	jwtSignPrivateKey *rsa.PrivateKey
)

const (
	jwtPrivateFileName = ".matticnote_jwt_private.pem"
	jwtPublicFileName  = ".matticnote_jwt_public.pem"
	authSchemeName     = "jwt"
	authHeaderName     = "Authorization"
	JWTAuthCookieName  = "jwt_auth"
)

func GenerateJWTSignKey(overwrite bool) error {
	var shouldGen = false
	_, err := os.Stat(jwtPrivateFileName)
	if err != nil {
		if os.IsNotExist(err) {
			shouldGen = true
		} else {
			return err
		}
	}
	_, err = os.Stat(jwtPublicFileName)
	if err != nil {
		if os.IsNotExist(err) {
			shouldGen = true
		} else {
			return err
		}
	}

	if shouldGen || overwrite {
		priKey, pubKey := misc.GenerateRSAKeypair(2048)
		err = ioutil.WriteFile(jwtPrivateFileName, priKey, fs.FileMode(0600))
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(jwtPublicFileName, pubKey, fs.FileMode(0600))
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadJWTSignKey() error {
	pubKeyByte, err := ioutil.ReadFile(jwtPublicFileName)
	if err != nil {
		return err
	}

	priKeyByte, err := ioutil.ReadFile(jwtPrivateFileName)
	if err != nil {
		return err
	}

	pubKeyPem, _ := pem.Decode(pubKeyByte)
	priKeyPem, _ := pem.Decode(priKeyByte)

	parsedPubKey, err := x509.ParsePKCS1PublicKey(pubKeyPem.Bytes)
	if err != nil {
		return err
	}

	parsedPriKey, err := x509.ParsePKCS1PrivateKey(priKeyPem.Bytes)
	if err != nil {
		return err
	}

	jwtSignPublicKey = parsedPubKey
	jwtSignPrivateKey = parsedPriKey

	return nil
}

func SignJWT(userUUID uuid.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userUUID.String()

	signed, err := token.SignedString(jwtSignPrivateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func RegisterFiberJWT(mode string) fiber.Handler {
	var authFailed = func(c *fiber.Ctx, err error) error {
		if c.Accepts("html") != "" {
			return c.Redirect(fmt.Sprintf("/account/login?next=%s", url.QueryEscape(c.OriginalURL())))
		} else {
			c.Status(401)
			return nil
		}
	}

	switch strings.ToLower(mode) {
	case "cookie":
		return jwtware.New(jwtware.Config{
			ErrorHandler:  authFailed,
			SigningKey:    jwtSignPublicKey,
			SigningMethod: "RS512",
			ContextKey:    "jwt_user",
			TokenLookup:   fmt.Sprintf("cookie:%s", JWTAuthCookieName),
		})
	default: // within header
		return jwtware.New(jwtware.Config{
			Filter: func(c *fiber.Ctx) bool {
				authHeaderSlice := strings.Split(c.Get("Authorization", ""), " ")
				if len(authHeaderSlice) > 0 && strings.TrimSpace(authHeaderSlice[0]) == authSchemeName {
					return true
				} else {
					return false
				}
			},
			ErrorHandler:  authFailed,
			SigningKey:    jwtSignPublicKey,
			SigningMethod: "RS512",
			ContextKey:    "jwt_user",
			TokenLookup:   fmt.Sprintf("header:%s", authHeaderName),
			AuthScheme:    authSchemeName,
		})
	}
}
