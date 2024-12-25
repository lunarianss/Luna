package route

import (
	"github.com/gin-gonic/gin"
	service "github.com/lunarianss/Luna/internal/api-server/application"
	"github.com/lunarianss/Luna/internal/api-server/config"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	controller "github.com/lunarianss/Luna/internal/api-server/interface/gin/v1/app"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
)

type AppRoutes struct{}

func (a *AppRoutes) Register(g *gin.Engine) error {
	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	// repos
	accountRepo := repo_impl.NewAccountRepoImpl(gormIns)
	tenantRepo := repo_impl.NewTenantRepoImpl(gormIns)
	appRepo := repo_impl.NewAppRepoImpl(gormIns)
	messageRepo := repo_impl.NewMessageRepoImpl(gormIns)
	providerRepo := repo_impl.NewProviderRepoImpl(gormIns)
	webAppRepo := repo_impl.NewWebAppRepoImpl(gormIns)
	modelProviderRepo := repo_impl.NewModelProviderRepoImpl(gormIns)
	providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelProviderRepo, "", nil)
	annotationRepo := repo_impl.NewAnnotationRepoImpl(gormIns)
	// domain
	providerDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, tenantRepo, providerConfigurationsManager)
	appDomain := appDomain.NewAppDomain(appRepo, webAppRepo, gormIns)
	accountDomain := accountDomain.NewAccountDomain(accountRepo, nil, nil, nil, tenantRepo)
	chatDomain := chatDomain.NewChatDomain(messageRepo, annotationRepo)

	// service
	appService := service.NewAppService(appDomain, providerDomain, accountDomain, gormIns, config)
	chatService := service.NewChatService(appDomain, providerDomain, accountDomain, chatDomain)

	appController := controller.NewAppController(appService, chatService)

	v1 := g.Group("/v1")
	modelProviderV1 := v1.Group("/console/api")
	modelProviderV1.Use(middleware.TokenAuthMiddleware())
	modelProviderV1.POST("/apps/:appID/site-enable", appController.UpdateAppSiteEnable)
	modelProviderV1.POST("/apps/:appID/api-enable", appController.UpdateAppApiEnable)
	modelProviderV1.POST("/rule-generate", appController.GeneratePrompt)
	modelProviderV1.POST("/apps/:appID/api-keys", appController.GenerateAppServiceToken)
	modelProviderV1.GET("/apps/:appID/api-keys", appController.ListAppServiceToken)

	return nil
}

func (r *AppRoutes) GetModule() string {
	return "console_app"
}
