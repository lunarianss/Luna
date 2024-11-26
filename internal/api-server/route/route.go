// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	authRoute "github.com/lunarianss/Luna/internal/api-server/route/console/auth"
	datasetRoute "github.com/lunarianss/Luna/internal/api-server/route/console/dataset"
	featureRoute "github.com/lunarianss/Luna/internal/api-server/route/console/feature"
	setupRoute "github.com/lunarianss/Luna/internal/api-server/route/console/setup"
	consoleWorkSpaceRoute "github.com/lunarianss/Luna/internal/api-server/route/console/workspace"
	chatAppRoute "github.com/lunarianss/Luna/internal/api-server/route/web/chat_app"
	"github.com/lunarianss/Luna/internal/pkg/server"
)

// Route unified registration portal
func init() {
	server.RegisterRoute(&blogRoutes{})

	// console/workspace
	server.RegisterRoute(&consoleWorkSpaceRoute.ModelProviderRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.ModelRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.AppRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.AccountRoute{})
	server.RegisterRoute(&consoleWorkSpaceRoute.WorkspaceRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.TagRoutes{})

	// console/dataset
	server.RegisterRoute(&datasetRoute.DatasetRoutes{})

	// console/auth
	server.RegisterRoute(&authRoute.AuthRoutes{})

	// console/setup
	server.RegisterRoute(&setupRoute.SetupRoutes{})

	// console/feature
	server.RegisterRoute(&featureRoute.FeatureRoutes{})
	server.RegisterRoute(&staticRoute{})

	// web
	server.RegisterRoute(&chatAppRoute.PassportRoutes{})
	server.RegisterRoute(&chatAppRoute.WebSiteRoutes{})
	server.RegisterRoute(&chatAppRoute.WebAppRoutes{})
	server.RegisterRoute(&chatAppRoute.WebChatRoutes{})
	server.RegisterRoute(&chatAppRoute.WebMessageRoutes{})
}
