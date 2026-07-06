package delivery

import "github.com/gin-gonic/gin"

type rootDelivery struct {
	demoDelivery        *demoDelivery
	articleDelivery     *articleDelivery
	authDelivery        *authDelivery
	userDelivery        *userDelivery
	chatGroupDelivery   *chatGroupDelivery
	chatMessageDelivery *chatMessageDelivery
	searchDelivery      *searchDelivery
	orderDelivery       *OrderDelivery
}

func NewRootDelivery(demoDelivery *demoDelivery, articleDelivery *articleDelivery, authDelivery *authDelivery, userDelivery *userDelivery, chatGroupDelivery *chatGroupDelivery, chatMessageDelivery *chatMessageDelivery, searchDelivery *searchDelivery, orderDelivery *OrderDelivery) *rootDelivery {
	return &rootDelivery{
		demoDelivery:        demoDelivery,
		articleDelivery:     articleDelivery,
		authDelivery:        authDelivery,
		userDelivery:        userDelivery,
		chatGroupDelivery:   chatGroupDelivery,
		chatMessageDelivery: chatMessageDelivery,
		searchDelivery:      searchDelivery,
		orderDelivery:       orderDelivery,
	}
}

func (r *rootDelivery) RegisterRouter(ginEngine *gin.Engine) {
	apiGroup := ginEngine.Group("api")
	{
		r.demoDelivery.RegisterRouter(apiGroup)
		r.articleDelivery.RegisterRouter(apiGroup)
		r.authDelivery.RegisterRouter(apiGroup)
		r.userDelivery.RegisterRouter(apiGroup)
		r.chatGroupDelivery.RegisterRouter(apiGroup)
		r.chatMessageDelivery.RegisterRouter(apiGroup)
		r.searchDelivery.RegisterRouter(apiGroup)
		r.orderDelivery.RegisterRouter(apiGroup)
	}
}
