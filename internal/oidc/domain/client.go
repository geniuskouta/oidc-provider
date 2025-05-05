package domain

type Client struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Name         string   `json:"name"`
	RedirectURIs []string `json:"redirect_uris"`
	CreatedAt    int64    `json:"created_at"`
}

type ClientRepository interface {
	Save(client *Client) error
	FindByID(clientID string) (*Client, error)
}
