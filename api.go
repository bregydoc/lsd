package lsd

import (
	"github.com/bregydoc/lsd/api"
	"github.com/rs/xid"
)


func (lsd *LSD) SendNotificationHTTP(p *api.NotificationPayload) (*api.NotificationResult, error) {
	id := xid.New().String()
	for _, to := range p.To {
		if err := lsd.emitNotification(&Notification{
			ID:          id,
			To:          to,
			Title:       p.Notification.Title,
			Body:        Markdown(p.Notification.Body),
			Options:     p.Notification.Options,
		}); err != nil {
			return nil, err
		}
	}

	return &api.NotificationResult{Ok: true, NotificationID: id}, nil
}

func (lsd *LSD) GenerateNewKeyPairHTTP(p *api.NewKeyPairPayload) (*api.KeyPairResult, error) {
	publicKey, privateKey, err := lsd.generateNewKeyPair()
	if err != nil {
		return nil, err
	}

	if err = lsd.saveKeyPair(p.UserID, publicKey, privateKey); err != nil {
		return nil, err
	}

	return &api.KeyPairResult{UserID: p.UserID, PublicKey: publicKey}, nil
}

func (lsd *LSD) GetKeyPairHTTP(p *api.KeyPairPayload) (*api.KeyPairResult, error) {
	public, _, err := lsd.getKeyPair(p.UserID)
	if err != nil {
		return nil, err
	}

	return &api.KeyPairResult{UserID: p.UserID, PublicKey: public}, nil
}
