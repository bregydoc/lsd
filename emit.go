package lsd

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

func (lsd *LSD) emitNotification(notification *Notification) error {
	notification.DeliveredAt = time.Now()

	if lsd.sessionsMap == nil {
		log.Error("sessions map not initialized")
		return errors.New("sessions map not initialized")
	}

	s, sessionExist := lsd.sessionsMap[notification.To]
	if !sessionExist {
		log.Error("user session not registered")
		return errors.New("user session not registered")
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	log.Info("before: ", string(payload))

	// privateKey, err := lsd.getPrivateKey(notification.To)
	// if err != nil {
	// 	return err
	// }
	//
	// payload, err = lsd.encryptNotification(privateKey, string(payload))
	// if err != nil {
	// 	return err
	// }

	token, err := lsd.getToken(notification.To)
	if err != nil {
		return err
	}

	payload, err = lsd.encryptNotification(token, payload)
	if err != nil {
		return err
	}

	log.Info("after: ", string(payload))

	payload64 := base64.StdEncoding.EncodeToString(payload)

	return s.Write([]byte(payload64))
}
