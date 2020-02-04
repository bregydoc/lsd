package lsd

import (
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
		if err := m.HandleRequest(c.Writer, c.Request); err != nil {
			panic(err)
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Info("connecting... ", s.Request.Host)
		lock.Lock()

		uri, err := url.Parse(s.Request.RequestURI)
		if err != nil {
			log.Error(err)
			return
		}

		values := uri.Query()

		userID := values.Get("userID")
		// publicKey := values.Get("publicKey")
		token := values.Get("token")

		decodedUserID, err := base64.StdEncoding.DecodeString(userID)
		if err != nil {
			log.Error(err)
			return
		}

		decodeToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			log.Error(err)
			return
		}

		userID = string(decodedUserID)
		token = string(decodeToken)

		log.Info("userID: ", userID)
		log.Info("token:  ", token)

		isMono := token == ""

		if lsd.secure && isMono {
			log.Error("lsd.secure && isMono")
			return // TODO: Send an error message
		}

		if isMono {
			lsd.sessionsMap[userID] = s
			// The session id of a user is the same id, in this case
			if err := lsd.registerUserSession(userID, userID); err != nil {
				log.Error(err)
				return // TODO: Send an error message
			}
		} else {
			// if err := lsd.ifPublicKeyMatchWithUserID(userID, p64); err != nil {
			// 	log.Error(err)
			// 	return

			if err := lsd.ifTokenMatchWithUserID(userID, token); err != nil {
				log.Error(err)
				return
			}

			lsd.sessionsMap[userID] = s
			// The session id of a user is the same id, in this case
			if err := lsd.registerUserSession(userID, userID); err != nil {
				log.Error(err)
				return // TODO: Send an error message
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
