package lsd

import (
	"bytes"
	"encoding/base64"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)

func (lsd *LSD) launchClientWSServer(addr ...string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	m := melody.New()
	lock := new(sync.Mutex)


	if lsd.sessionsMap == nil {
		lsd.sessionsMap = map[string]*melody.Session{}
	}

	log.Info("registering lsd router for ws")
	r.GET("/ws", func(c *gin.Context) {
		log.Info("/ws ", c.Request.Host)
		if err := m.HandleRequest(c.Writer, c.Request); err != nil {
			panic(err)
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Info("connecting... ", s.Request.Host)
		lock.Lock()

		log.Info(s.Request.RequestURI)
		values, err := url.ParseQuery(s.Request.RequestURI)
		if err != nil {
			log.Error(err)
			return
		}
		// protocolHeader := s.Request.Header.Get("Sec-WebSocket-Protocol")
		userIDFromQuery := values.Get("userID")
		publicKeyFromQuery := values.Get("publicKey")

		log.Info(userIDFromQuery, publicKeyFromQuery)
		// chunks :=  strings.Split(protocolHeader, ",")

		userID := userIDFromQuery
		publicKey := publicKeyFromQuery
		//
		// if len(chunks) > 1 {
		// 	publicKey = chunks[1]
		// }

		log.Info("userID, publicKey ", userID, publicKey)

		isMono := publicKey == ""

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
			private, err := lsd.getPrivateKey(userID)
			log.Info("private key x: ", string(private))
			if err != nil {
				log.Error(err)
				return
			}

			public, err := lsd.publicKeyBytesFromPrivateKeyBytes(private)
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
		if err := s.Write([]byte("PONG")); err != nil {
			log.Error(err)
			return
		}

		lock.Unlock()
	})

	// m.HandleMessage(func(s *melody.Session, msg []byte) {
	//
	// })

	return r.Run(addr...)
}