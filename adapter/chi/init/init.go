package init

import (
	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/adapter/chi"
)

func init() {
	oas.RegisterAdapter("chi", oas_chi.NewAdapter())
}
