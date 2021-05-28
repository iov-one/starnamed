package pkg

import (
	"github.com/iov-one/starnamed/app"
)

// ModuleCdc instantiates a new codec for the domain module
var ModuleCdc = app.MakeEncodingConfig()
