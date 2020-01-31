package api


type NotificationPayload struct {
	To                   []string
	Notification         Notification
}

type NotificationResult struct {
	Ok                   bool
	NotificationID       string
}

type Notification struct {
	Title                string
	Body                 string
	Options              []string
}

type NewKeyPairPayload struct {
	UserID               string
}

type KeyPairResult struct {
	UserID               string
	PublicKey            []byte
}

type KeyPairPayload struct {
	UserID               string
}

type ServiceHTTP interface {
	SendNotificationHTTP(*NotificationPayload) (*NotificationResult, error)
	GenerateNewKeyPairHTTP(*NewKeyPairPayload) (*KeyPairResult, error)
	GetKeyPairHTTP(*KeyPairPayload) (*KeyPairResult, error)
}

