package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/domain_service"
	repo_impl "github.com/lunarianss/Luna/internal/api-server/_repo"
	"github.com/lunarianss/Luna/internal/api-server/config"
	controller "github.com/lunarianss/Luna/internal/api-server/controller/gin/v1/web/chat_app/chat"
	"github.com/lunarianss/Luna/internal/api-server/dao"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/app"
	appRunningDomain "github.com/lunarianss/Luna/internal/api-server/domain/app_running"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat"
	"github.com/lunarianss/Luna/internal/api-server/middleware"
	"github.com/lunarianss/Luna/internal/api-server/service"
	"github.com/lunarianss/Luna/internal/pkg/email"
	"github.com/lunarianss/Luna/internal/pkg/mysql"
	"github.com/lunarianss/Luna/internal/pkg/redis"
)

type WebChatRoutes struct{}

func (a *WebChatRoutes) Register(g *gin.Engine) error {

	gormIns, err := mysql.GetMySQLIns(nil)

	if err != nil {
		return err
	}

	redisIns, err := redis.GetRedisIns(nil)

	if err != nil {
		return err
	}

	email, err := email.GetEmailSMTPIns(nil)

	if err != nil {
		return err
	}

	// config
	config, err := config.GetLunaRuntimeConfig()

	if err != nil {
		return err
	}

	// dao
	appDao := dao.NewAppDao(gormIns)
	appRunningDao := dao.NewAppRunningDao(gormIns)
	accountDao := dao.NewAccountDao(gormIns)
	tenantDao := dao.NewTenantDao(gormIns)
	messageDao := dao.NewMessageDao(gormIns)

	providerRepo := repo_impl.NewProviderRepoImpl(gormIns)
	modelProviderRepo := repo_impl.NewModelProviderRepoImpl(gormIns)
	providerConfigurationsManager := domain_service.NewProviderConfigurationsManager(providerRepo, modelProviderRepo, "", nil)

	// domain
	appDomain := domain.NewAppDomain(appDao, appRunningDao, messageDao)
	appRunningDomain := appRunningDomain.NewAppRunningDomain(appRunningDao)
	accountDomain := accountDomain.NewAccountDomain(accountDao, redisIns, config, email, tenantDao)
	// providerDomain := providerDomain.NewModelProviderDomain(providerDao, modelDao)
	chatDomain := chatDomain.NewChatDomain(messageDao)

	// domain
	providerDomain := domain_service.NewProviderDomain(providerRepo, modelProviderRepo, providerConfigurationsManager)
	webChatService := service.NewWebChatService(appRunningDomain, accountDomain, appDomain, config, providerDomain, chatDomain)

	webSiteController := controller.NewWebChatController(webChatService)
	v1 := g.Group("/v1")
	authV1 := v1.Group("/api")
	authV1.Use(middleware.WebTokenAuthMiddleware())
	authV1.POST("/chat-messages", webSiteController.Chat)
	return nil
}

func (r *WebChatRoutes) GetModule() string {
	return "web_chat"
}
