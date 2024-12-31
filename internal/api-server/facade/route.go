// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import (
	consoleAppRoute "github.com/lunarianss/Luna/internal/api-server/facade/console/app"
	authRoute "github.com/lunarianss/Luna/internal/api-server/facade/console/auth"
	datasetRoute "github.com/lunarianss/Luna/internal/api-server/facade/console/dataset"
	featureRoute "github.com/lunarianss/Luna/internal/api-server/facade/console/feature"
	setupRoute "github.com/lunarianss/Luna/internal/api-server/facade/console/setup"
	consoleWorkSpaceRoute "github.com/lunarianss/Luna/internal/api-server/facade/console/workspace"
	chatAppRoute "github.com/lunarianss/Luna/internal/api-server/facade/web/chat_app"
	"github.com/lunarianss/Luna/internal/infrastructure/server"

	serviceChatRoute "github.com/lunarianss/Luna/internal/api-server/facade/service/chat_app"
)

// Route unified registration portal
func init() {
	// console/workspace
	server.RegisterRoute(&consoleWorkSpaceRoute.ModelProviderRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.ModelRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.AppRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.AccountRoute{})
	server.RegisterRoute(&consoleWorkSpaceRoute.WorkspaceRoutes{})
	server.RegisterRoute(&consoleWorkSpaceRoute.TagRoutes{})

	// console/app
	server.RegisterRoute(&consoleAppRoute.ChatRoutes{})
	server.RegisterRoute(&consoleAppRoute.StatisticRoutes{})
	server.RegisterRoute(&consoleAppRoute.AppRoutes{})
	server.RegisterRoute(&consoleAppRoute.AnnotationRoutes{})
	server.RegisterRoute(&consoleAppRoute.ToolRoutes{})

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

	// service
	server.RegisterRoute(&serviceChatRoute.ServiceChatRoutes{})
}
