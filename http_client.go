package nuwa

import "github.com/parnurzeal/gorequest"

var superAgent *gorequest.SuperAgent

func HttpClient() *gorequest.SuperAgent {
	if superAgent == nil {
		superAgent = gorequest.New()
	}
	return superAgent
}
