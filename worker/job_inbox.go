package worker

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/go-fed/httpsig"
	"github.com/gocraft/work"
	"github.com/piprate/json-gold/ld"
	"log"
)

var (
	jldProc    = ld.NewJsonLdProcessor()
	jldOptions = ld.NewJsonLdOptions("")
	jldDoc     = map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
	}
)

func (c *Context) ProcessInbox(j *work.Job) error {
	dataRaw, ok := j.Args["data"]
	if !ok {
		log.Println("err: 'data' argument is not found. ignore.")
		return nil
	}
	data, ok := dataRaw.(map[string]interface{})
	if !ok {
		return errors.New("cannot convert to map[string]interface{}")
	}

	doc, err := jldProc.Compact(data, jldDoc, jldOptions)
	if err != nil {
		// parse failed
		return nil
	}
	if len(doc) == 0 {
		// no length
		return nil
	}

	actorRaw, ok := doc["actor"]
	if ok {
		actor, ok := actorRaw.(string)
		if !ok {
			// unknown actor type
			return nil
		}
		_, err := user.GetRemoteUserFromApID(actor)
		if err == user.ErrNoSuchUser {
			_, err := user.RegisterRemoteUser(actor)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		verifierRaw, ok := j.Args["signature"]
		if ok {
			verifier := verifierRaw.(httpsig.Verifier)

			userPK, err := user.GetUserPublicKeyFromKeyId(verifier.KeyId())
			if err != nil {
				log.Println("err: signature missing. ignore.")
				return nil
			}
			err = verifier.Verify(userPK, httpsig.RSA_SHA256)
			if err != nil {
				log.Println("err: invalid HTTP signature. ignore.")
				return nil
			}
		} else {
			log.Println("err: HTTP signature is not defined. ignore.")
			return nil
		}
	} else {
		log.Println("err: unknown actor. ignore.")
		return nil
	}

	apTypeRaw, ok := doc["type"]
	if !ok {
		log.Println("err: type is not defined. ignore.")
		return nil
	}
	apType, ok := apTypeRaw.(string)
	if !ok {
		log.Println("err: activity type is not string. ignore.")
		return nil
	}

	switch apType {
	case "Create":
		log.Println("Create activity")
	default:
		log.Println("err: unknown activity. ignore.")
	}

	return nil
}
