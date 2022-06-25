package rgvalidator

import "rgo/core/rgrequest"

type Validator interface {
	CheckParam(*rgrequest.Client)
}
