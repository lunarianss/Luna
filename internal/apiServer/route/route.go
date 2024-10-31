// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import "github.com/Ryan-eng-del/hurricane/internal/pkg/server"

// Route unified registration portal
func init() {
	server.RegisterRoute(&blogRoutes{})
	server.RegisterRoute(&staticRoute{})
}
