package swagger

import (
	"go-backend/ent"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

func auth(reflector *openapi3.Reflector) error {
	{
		operationContext, err := reflector.NewOperationContext(http.MethodPost, "/api/auth/login")
		if err != nil {
			return err
		}
		operationContext.SetTags("Auth")
		operationContext.SetSummary("Login")
		operationContext.SetDescription("Đăng nhập với email/pass")

		operationContext.AddReqStructure(new(dto.AuthLoginReq))
		operationContext.AddRespStructure(new(response.SuccessFormat[bool]))

		reflector.AddOperation(operationContext)
	}

	{
		operationContext, err := reflector.NewOperationContext(http.MethodGet, "/api/auth/get-info")
		if err != nil {
			return err
		}
		operationContext.SetTags("Auth")
		operationContext.SetSummary("Lay thong tin user")
		operationContext.SetDescription("description")

		operationContext.AddRespStructure(new(*ent.Users))

		reflector.AddOperation(operationContext)
	}

	{
		operationContext, err := reflector.NewOperationContext(http.MethodPost, "/api/auth/refresh-token")
		if err != nil {
			return err
		}
		operationContext.SetTags("Auth")
		operationContext.SetSummary("Lam moi access token")
		operationContext.SetDescription("description")

		operationContext.AddRespStructure(new(response.SuccessFormat[bool]))

		reflector.AddOperation(operationContext)
	}
	return nil
}
