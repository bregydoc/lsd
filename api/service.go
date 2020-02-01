package api

type NotificationPayload struct {
	To           []string
	Notification Notification
}

type NotificationResult struct {
	Ok            bool
	Notifications map[string]string
}

type Notification struct {
	Title   string
	Body    string
	Options []string
}

type NewKeyPairPayload struct {
	UserID string
}

type KeyPairResult struct {
	UserID    string
	PublicKey []byte
}

type KeyPairPayload struct {
	UserID string
}

type NewTokenPayload struct {
	UserID string
}

type TokenResult struct {
	UserID string
	Token  string
}

type TokenPayload struct {
	UserID string
}

type ServiceHTTP interface {
	SendNotificationHTTP(*NotificationPayload) (*NotificationResult, error)
	// GenerateNewKeyPairHTTP(*NewKeyPairPayload) (*KeyPairResult, error)
	GenerateNewTokenForUserHTTP(*NewTokenPayload) (*TokenResult, error)
	// GetKeyPairHTTP(*KeyPairPayload) (*KeyPairResult, error)
	GetTokenHTTP(*TokenPayload) (*TokenResult, error)
}
