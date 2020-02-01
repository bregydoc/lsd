package lsd

import (
	"github.com/bregydoc/lsd/api"
	"github.com/rs/xid"
)

func (lsd *LSD) SendNotificationHTTP(p *api.NotificationPayload) (*api.NotificationResult, error) {
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

	return &api.NotificationResult{Ok: true, Notifications: notifications}, nil
}

// func (lsd *LSD) GenerateNewKeyPairHTTP(p *api.NewKeyPairPayload) (*api.KeyPairResult, error) {
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
// 	return &api.KeyPairResult{UserID: p.UserID, PublicKey: publicKey}, nil
// }

func (lsd *LSD) GenerateNewTokenForUserHTTP(p *api.NewTokenPayload) (*api.TokenResult, error) {
	token, err := lsd.generateNewToken()
	if err != nil {
		return nil, err
	}

	if err = lsd.saveToken(p.UserID, token); err != nil {
		return nil, err
	}

	return &api.TokenResult{
		UserID: p.UserID,
		Token:  token,
	}, nil
}

// func (lsd *LSD) GetKeyPairHTTP(p *api.KeyPairPayload) (*api.KeyPairResult, error) {
// 	privateKey, err := lsd.getPrivateKey(p.UserID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	publicKey, err := lsd.publicKeyBytesFromPrivateKeyBytes(privateKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &api.KeyPairResult{UserID: p.UserID, PublicKey: publicKey}, nil
// }

func (lsd *LSD) GetTokenHTTP(p *api.TokenPayload) (*api.TokenResult, error) {
	token, err := lsd.getToken(p.UserID)
	if err != nil {
		return nil, err
	}

	return &api.TokenResult{UserID: p.UserID, Token: token}, nil
}
