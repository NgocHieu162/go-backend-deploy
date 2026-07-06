package elastic

import (
	"context"
	"fmt"
	"go-backend/internal/common/env"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/v9"
)

type Elastic struct {
	EsClient          *elasticsearch.TypedClient
	articleRepository interfaces.ArticleRepository
	userRepository    interfaces.UserRepository
}

func NewElastic(env *env.Env, articleRepository interfaces.ArticleRepository, userRepository interfaces.UserRepository) *Elastic {
	esClient, err := elasticsearch.NewTyped(
		elasticsearch.WithAddresses(env.ElasticAddr),
		elasticsearch.WithBasicAuth(env.ElasticUser, env.ElasticPass),
		elasticsearch.WithCertificateFingerprint(env.ElasticCertFingerprint),
	)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	info, err := esClient.Info().Do(context.Background())
	if err != nil {
		log.Fatal("Elastic Connect Error: ", err)
		return nil
	}

	fmt.Println("[ELASTIC] Connect to Elastic Successfully", info.ClusterName)

	return &Elastic{
		EsClient:          esClient,
		articleRepository: articleRepository,
		userRepository:    userRepository,
	}
}

func (e *Elastic) InitArticle() {
	ctx := context.Background()
	articles, err := e.articleRepository.GetAll(ctx,
		pagination.Query{
			Page:     1,
			PageSize: 99999,
		},
		dto.ArticleFindAllFilters{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, article := range articles {
		e.EsClient.Index("articles").
			Id(strconv.Itoa(article.ID)).
			Document(article).
			Do(ctx)
	}
}

func (e *Elastic) InitUser() {
	ctx := context.Background()
	users, err := e.userRepository.GetAll(ctx,
		pagination.Query{
			Page:     1,
			PageSize: 99999,
		},
		dto.UserFindAllFilters{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, user := range users {
		e.EsClient.Index("users").
			Id(strconv.Itoa(user.ID)).
			Document(user).
			Do(ctx)
	}
}
