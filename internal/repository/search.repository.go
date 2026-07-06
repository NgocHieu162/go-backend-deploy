package repository

import (
	"context"
	"go-backend/internal/common/elastic"
	"go-backend/internal/interfaces"

	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/operator"
)

type SearchRepository struct {
	elastic *elastic.Elastic
}

func NewSearchRepository(elastic *elastic.Elastic) interfaces.SearchRepository {
	return &SearchRepository{
		elastic: elastic,
	}
}

// FindAll implements [repository.SearchRepository].
func (a *SearchRepository) FindAll(ctx context.Context, textSearch string) (any, error) {
	return a.elastic.EsClient.Search().
		Index("articles,users").
		Request(
			&search.Request{
				Query: &types.Query{
					MultiMatch: &types.MultiMatchQuery{
						Query: textSearch,
						Fields: []string{
							"title",
							"content",
							"email",
							"fullName",
						},
						// Với nhiều từ khóa, chỉ cân khớp một phần cũng được
						Operator: &operator.Or,

						// Cho phép user gõ sai nhẹ vẫn tìm
						Fuzziness: "AUTO",

						// Document nên khớp khoảng 60%
						MinimumShouldMatch: "60%",
					},
				},
			},
		).Do(ctx)
}
