package rgvalidator

import "github.com/jackylee92/rgo/core/rgrequest"

type Validator interface {
	CheckParam(*rgrequest.Client)
}
