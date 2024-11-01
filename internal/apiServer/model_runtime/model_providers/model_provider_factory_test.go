package model_providers

import (
	"encoding/json"
	"testing"

	"github.com/lunarianss/Hurricane/pkg/errors"
	"github.com/lunarianss/Hurricane/pkg/log"
)

func TestModelProviderFactory(t *testing.T) {
	mf := ModelProviderFactory{}
	log.NewWithOptions(log.WithDebugMode())
	providers, err := mf.GetProvidersFromDir()

	if err != nil {
		if coder, ok := err.(errors.Coder); ok {
			t.Logf("%#+v", coder)
		} else {
			t.Log(err.Error())
		}
		return
	}

	c, _ := json.MarshalIndent(providers, "", " ")
	t.Logf("providers: %+v", string(c))
}
