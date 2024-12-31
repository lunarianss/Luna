package tools

import (
	"testing"

	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func TestListBuiltInProvider(t *testing.T) {
	log.NewWithOptions(log.WithDebugMode())

	d, err := NewToolManager().ListBuiltInProviders()

	if err != nil {
		log.Errorf("%#+v", err)
	}

	util.LogCompleteInfo(d)
}
