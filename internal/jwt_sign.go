package internal

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/form3tech-oss/jwt-go"
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
	AuthSchemeName     = "jwt"
	AuthHeaderName     = "Authorization"
	JWTAuthCookieName  = "jwt_auth"
	jwtSignMethod      = "RS512"
	JWTContextKey      = "jwt_user"
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
	if pubKeyPem.Type != "PUBLIC KEY" {
		return errors.New("this is not public key")
	}
	priKeyPem, _ := pem.Decode(priKeyByte)
	if priKeyPem.Type != "PRIVATE KEY" {
		return errors.New("this is not private key")
	}

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

func VerifyRSASign() error {
	testHash := sha256.New()
	_, err := testHash.Write([]byte(misc.GenToken(16)))
	if err != nil {
		panic(err)
	}
	testHashSum := testHash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, jwtSignPrivateKey, crypto.SHA256, testHashSum, nil)
	if err != nil {
		panic(err)
	}

	err = rsa.VerifyPSS(jwtSignPublicKey, crypto.SHA256, testHashSum, signature, nil)
	if err != nil {
		if err == rsa.ErrVerification {
			return errors.New("the key pair does not match. If the problem persists, try deleting the key file")
		} else {
			panic(err)
		}
	}

	// Verify success
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

func RegisterFiberJWT(mode string, mustLogin bool) fiber.Handler {
	var authFailed = func() fiber.ErrorHandler {
		if mustLogin {
			return func(c *fiber.Ctx, _ error) error {
				if c.Accepts("html") != "" {
					return c.Redirect(fmt.Sprintf("/account/login?next=%s", url.QueryEscape(c.OriginalURL())))
				} else {
					c.Status(401)
					return nil
				}
			}
		} else {
			return func(c *fiber.Ctx, _ error) error {
				return c.Next()
			}
		}
	}()

	switch strings.ToLower(mode) {
	case "cookie":
		return jwtware.New(jwtware.Config{
			ErrorHandler:  authFailed,
			SigningKey:    jwtSignPublicKey,
			SigningMethod: jwtSignMethod,
			ContextKey:    JWTContextKey,
			TokenLookup:   fmt.Sprintf("cookie:%s", JWTAuthCookieName),
		})
	default: // within header
		return jwtware.New(jwtware.Config{
			ErrorHandler:  authFailed,
			SigningKey:    jwtSignPublicKey,
			SigningMethod: jwtSignMethod,
			ContextKey:    JWTContextKey,
			TokenLookup:   fmt.Sprintf("header:%s", AuthHeaderName),
			AuthScheme:    AuthSchemeName,
		})
	}
}
