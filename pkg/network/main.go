package network

import "github.com/toastcheng/dicomweb-go/dicomweb"

type Client struct {
	client *dicomweb.Client
}

func NewClient(qidoURL, wadoURL, stowURL string) Client {
	return Client{client: dicomweb.NewClient(dicomweb.ClientOption{
		QIDOEndpoint: qidoURL,
		WADOEndpoint: wadoURL,
		STOWEndpoint: stowURL,
	})}
}
