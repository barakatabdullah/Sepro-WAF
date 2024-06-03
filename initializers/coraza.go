package initializers

import (
	"github.com/corazawaf/coraza/v3"
)

func CreateWaf() (coraza.WAF, error) {

	wafConfig := coraza.NewWAFConfig()

	waf, err := coraza.NewWAF(wafConfig)
	if err != nil {
		return nil, err
	}

	return waf, nil
}
