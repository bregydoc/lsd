package lsd

import (
	"bytes"
	"encoding/base64"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)



func (lsd *LSD) launchClientWSServer(addr ...string) error {
	r := gin.New()
	m := melody.New()
	lock := new(sync.Mutex)

	if lsd.sessionsMap == nil {
		lsd.sessionsMap = map[string]*melody.Session{}
	}

	r.GET("/ws", func(c *gin.Context) {
		if err := m.HandleRequest(c.Writer, c.Request); err != nil {
			panic(err)
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		lock.Lock()
		protocolHeader := s.Request.Header.Get("Sec-WebSocket-Protocol")
		chunks :=  strings.Split(protocolHeader, ",")

		userID := chunks[0]
		publicKey := chunks[0]

		if len(chunks) > 1 {
			publicKey = chunks[1]
		}

		isMono := userID == publicKey

		if lsd.secure && isMono {
			log.Error("lsd.secure && isMono")
			return // TODO: Send an error message
		}

		if isMono {
			lsd.sessionsMap[userID] = s
			// The session id of a user is the same id, in this case
			if err := lsd.registerUserSession(userID, userID); err != nil {
				log.Error(err)
				return  // TODO: Send an error message
			}
		} else {
			public, _, err := lsd.getKeyPair(userID)
			if err != nil {
				log.Error(err)
				return
			}

			requestPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
			if err != nil {
				log.Error(err)
				return
			}

			if !bytes.Equal(public, requestPublicKey) {
				log.Error("invalid public key, it isn't equal to our saved key")
				return
			}

			lsd.sessionsMap[userID] = s
			// The session id of a user is the same id, in this case
			if err := lsd.registerUserSession(userID, userID); err != nil {
				log.Error(err)
				return  // TODO: Send an error message
			}
		}

		lock.Unlock()
	})

	// m.HandleMessage(func(s *melody.Session, msg []byte) {
	//
	// })

	if err := r.Run(addr...); err != nil {
		log.Error(err)
		return err
	}

	return nil
}