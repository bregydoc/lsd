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

	notifications := map[string]string{}
	for _, to := range p.To {
		id := xid.New().String()
		if err := lsd.emitNotification(&Notification{
			ID:      id,
			To:      to,
			Title:   p.Notification.Title,
			Body:    Markdown(p.Notification.Body),
			Options: p.Notification.Options,
		}); err != nil {
			return nil, err
		}
		notifications[to] = id
	}

	return &proto.NotificationResult{Ok: true, Notifications: notifications}, nil
}
//
// func (lsd *LSD) GenerateNewKeyPair(c context.Context, p *proto.NewKeyPairPayload) (*proto.KeyPairResult, error) {
// 	privateKey, err := lsd.generateNewKeyPair()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if err = lsd.savePrivateKey(p.UserID, privateKey); err != nil {
// 		return nil, err
// 	}
//
// 	publicKey, err := lsd.publicKeyBytesFromPrivateKeyBytes(privateKey)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &proto.KeyPairResult{UserID: p.UserID, PublicKey: publicKey}, nil
// }

func (lsd *LSD) GenerateNewTokenForUser(c context.Context, p *proto.NewTokenPayload) (*proto.TokenResult, error) {
	token, err := lsd.generateNewToken()
	if err != nil {
		return nil, err
	}

	if err = lsd.saveToken(p.UserID, token); err != nil {
		return nil, err
	}

	return &proto.TokenResult{
		UserID:               p.UserID,
		Token:                token,
	}, nil
}

// func (lsd *LSD) GetKeyPair(c context.Context, p *proto.KeyPairPayload) (*proto.KeyPairResult, error) {
// 	privateKey, err := lsd.getPrivateKey(p.UserID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	publicKey, err := lsd.publicKeyBytesFromPrivateKeyBytes(privateKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &proto.KeyPairResult{UserID: p.UserID, PublicKey: publicKey}, nil
// }

func (lsd *LSD) GetToken(c context.Context, p *proto.TokenPayload) (*proto.TokenResult, error) {
	token, err := lsd.getToken(p.UserID)
	if err != nil {
		return nil, err
	}

	return &proto.TokenResult{UserID: p.UserID, Token: token}, nil
}
