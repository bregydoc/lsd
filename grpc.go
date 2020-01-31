package lsd

import (
	"context"
	"errors"

	proto "github.com/bregydoc/lsd/proto"
	"github.com/rs/xid"
)



func (lsd *LSD) SendNotification(c context.Context, p *proto.NotificationPayload) (*proto.NotificationResult, error) {
	if p.Notification == nil {
		return nil, errors.New("invalid notification")
	}

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

	return &proto.NotificationResult{Ok: true, NotificationID: id}, nil
}

func (lsd *LSD) GenerateNewKeyPair(c context.Context, p *proto.NewKeyPairPayload) (*proto.KeyPairResult, error) {
	publicKey, privateKey, err := lsd.generateNewKeyPair()
	if err != nil {
		return nil, err
	}

	if err = lsd.savePrivateKey(p.UserID, publicKey, privateKey); err != nil {
		return nil, err
	}

	return &proto.KeyPairResult{UserID: p.UserID, PublicKey: publicKey}, nil
}

func (lsd *LSD) GetKeyPair(c context.Context, p *proto.KeyPairPayload) (*proto.KeyPairResult, error) {
	public, _, err := lsd.getPrivateKey(p.UserID)
	if err != nil {
		return nil, err
	}

	return &proto.KeyPairResult{UserID: p.UserID, PublicKey: public}, nil
}

