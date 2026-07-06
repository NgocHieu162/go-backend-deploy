package dependency

import (
	"go-backend/ent"
	"go-backend/internal/common/cache"
	"go-backend/internal/common/elastic"
	"go-backend/internal/common/env"
	"go-backend/internal/common/middlewares"
	"go-backend/internal/common/rabbitmq"
	"go-backend/internal/common/socket"
	storage_impl "go-backend/internal/common/storage/impl"
	"go-backend/internal/delivery"
	"go-backend/internal/handler"
	"go-backend/internal/repository"
	"go-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func Injection(ginEngine *gin.Engine, entClient *ent.Client, env *env.Env, AllowOrigins []string, rabbitmq *rabbitmq.RabbitMQ) {
	tokenUsecase := usecase.NewTokenUsecase(env)

	cache := cache.NewCache(env)

	localFileStorage := storage_impl.NewLocalFileStorage("public")
	cloudinaryStorage := storage_impl.NewCloudinaryStorage(env)

	articleRepository := repository.NewArticleRepository(entClient)
	userRepository := repository.NewUserRepository(entClient)
	chatGroupRepository := repository.NewChatGroupRepository(entClient)
	chatMessageRepository := repository.NewChatMessageRepository(entClient)
	chatGroupMembersRepository := repository.NewChatGroupMembersRepository(entClient)
	unitOfWorkRepository := repository.NewUnitOfWorkRepository(entClient)

	elastic := elastic.NewElastic(env, articleRepository, userRepository)
	elastic.InitArticle()
	elastic.InitUser()
	searchRepository := repository.NewSearchRepository(elastic)

	chatUsecase := usecase.NewChatUsecase(tokenUsecase, userRepository, chatGroupRepository, chatGroupMembersRepository, chatMessageRepository, unitOfWorkRepository)
	chatHandler := handler.NewChatHandler(chatUsecase)
	socket := socket.NewSocket(chatHandler)
	socket.Start(ginEngine, AllowOrigins)

	chatGroupUsecase := usecase.NewChatGroupUsecase(chatGroupRepository)
	chatMessageUsecase := usecase.NewChatMessageUsecase(chatMessageRepository)
	searchUsecase := usecase.NewSearchUsecase(searchRepository)
	articleUsecase := usecase.NewArticleUsecase(articleRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, localFileStorage, cloudinaryStorage)
	authUsecase := usecase.NewAuthUsecase(userRepository, tokenUsecase, env)
	demoUsecase := usecase.NewDemoUsecase()

	authMiddleware := middlewares.NewAuthMiddleware(tokenUsecase, userRepository)

	articleHandler := handler.NewArticleHandler(articleUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	authHandler := handler.NewAuthHandler(authUsecase, env)
	demoHandler := handler.NewDemoHandler(demoUsecase)
	chatGroupHandler := handler.NewChatGroupHandler(chatGroupUsecase)
	chatMessageHandler := handler.NewChatMessageHandler(chatMessageUsecase)
	searchHandler := handler.NewSearchHandler(searchUsecase)

	articleDelivery := delivery.NewArticleDelivery(articleHandler, cache, authMiddleware)
	userDelivery := delivery.NewUserDelivery(userHandler, authMiddleware)
	authDelivery := delivery.NewAuthDelivery(authHandler, authMiddleware)
	demoDelivery := delivery.NewDemoDelivery(demoHandler)
	chatGroupDelivery := delivery.NewChatGroupDelivery(chatGroupHandler)
	chatMessageDelivery := delivery.NewChatMessageDelivery(chatMessageHandler)
	searchDelivery := delivery.NewSearchDelivery(searchHandler)

	orderUsecase := usecase.NewOrderUsecase(rabbitmq)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	orderDelivery := delivery.NewOrderDelivery(orderHandler, authMiddleware)

	rootDelivery := delivery.NewRootDelivery(demoDelivery, articleDelivery, authDelivery, userDelivery, chatGroupDelivery, chatMessageDelivery, searchDelivery, orderDelivery)
	rootDelivery.RegisterRouter(ginEngine)
}
