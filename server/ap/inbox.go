package ap

import (
	"bytes"
	"encoding/json"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/worker"
	"github.com/gocraft/work"
	"github.com/gofiber/fiber/v2"
	"github.com/piprate/json-gold/ld"
	"log"
	"net/http"
	"strings"
)

var (
	jldProc    = ld.NewJsonLdProcessor()
	jldOptions = ld.NewJsonLdOptions("")
	jldDoc     = map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
	}
)

func inboxGet(c *fiber.Ctx) error {
	c.Status(fiber.StatusMethodNotAllowed)
	return nil
}

func inboxPost(w http.ResponseWriter, r *http.Request) {
	var err error
	if !strings.HasPrefix(
		strings.ToLower(r.Header.Get("Content-Type")),
		"application/activity+json") &&
		!strings.HasPrefix(
			strings.ToLower(r.Header.Get("Content-Type")),
			"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"") {
		// Invalid header. return error.
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bufBody := new(bytes.Buffer)
	_, err = bufBody.ReadFrom(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var inboxData map[string]interface{}
	if err := json.Unmarshal(bufBody.Bytes(), &inboxData); err != nil {
		// Not json. return error.
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	doc, err := jldProc.Compact(inboxData, jldDoc, jldOptions)
	if err != nil {
		log.Println("err: json-ld parse failed. ignore.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(doc) == 0 {
		log.Println("err: json-ld parsed, but nobody attributes. ignore.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	actorRaw, ok := doc["actor"]
	if ok {
		actor, ok := actorRaw.(string)
		if !ok {
			log.Println("err: unknown actor value type. ignore.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err := internal.GetRemoteUserFromApID(actor)
		if err == internal.ErrNoSuchUser {
			_, err := internal.RegisterRemoteUser(actor)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// ↓一時的にコメント化させてるだけなので、終わったら戻す
	//verifier, err := httpsig.NewVerifier(r)
	//if err != nil {
	//	// not found HTTP signature. return error.
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	//
	//userPK, err := internal.GetUserPublicKey(verifier.KeyId())
	//if err != nil {
	//	// Signature missing. ignore.
	//	return
	//}
	//err = verifier.Verify(userPK, httpsig.RSA_SHA256)
	//if err != nil {
	//	// Invalid HTTP Signature. return error.
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	_, err = worker.Enqueue.Enqueue(worker.JobInboxProcess, work.Q{
		"doc": doc,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
