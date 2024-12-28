// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_providers

import (
	"testing"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
)

func TestModelProviderFactory(t *testing.T) {
	mf := ModelProviderFactory{}
	log.NewWithOptions(log.WithDebugMode())
	providers, orderedProviders, err := mf.GetProvidersFromDir()

	if err != nil {
		if coder, ok := err.(errors.Coder); ok {
			t.Logf("%#+v", coder)
		} else {
			t.Log(err.Error())
		}
		return
	}

	// c, _ := json.MarshalIndent(providers, "", " ")
	t.Logf(
		"len providers : %d, orderedProviders %+v, the first three provider names are %s | %s | %s",
		len(providers),
		orderedProviders,
		providers[0].Provider,
		providers[1].Provider,
		providers[2].Provider,
	)
}
