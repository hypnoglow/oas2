package init

import (
	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/adapter/gorilla"
)

func init() {
	oas.RegisterAdapter("gorilla", oas_gorilla.NewAdapter())
}
