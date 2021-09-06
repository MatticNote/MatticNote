package graphql

import (
	"context"
	"github.com/MatticNote/MatticNote/internal/auth"
	"github.com/MatticNote/MatticNote/internal/ist"
	"github.com/MatticNote/MatticNote/server/api/graphql/common"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

type graphqlForm struct {
	Query     string                 `form:"query"`
	Operation string                 `form:"operation"`
	Variables map[string]interface{} `form:"variables"`
}

func GQLEndpoint(c *fiber.Ctx) error {
	formData := new(graphqlForm)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	usrContext := context.Background()

	if currentUsr, ok := c.Locals(auth.LoginUserLocal).(*ist.LocalUserStruct); ok {
		usrContext = context.WithValue(usrContext, common.ContextCurrentUser, currentUsr.Uuid)
	}
	if method, ok := c.Locals(auth.AuthorizeMethodLocal).(auth.AuthorizeMethod); ok {
		usrContext = context.WithValue(usrContext, common.ContextAuthorizeMethod, method)
	}
	if oauthToken, ok := c.Locals(auth.OAuthTokenLocal).(string); ok {
		usrContext = context.WithValue(usrContext, common.ContextOAuthToken, oauthToken)
	}
	c.SetUserContext(usrContext)

	result := graphql.Do(graphql.Params{
		Context:        c.UserContext(),
		Schema:         setupSchema(),
		RequestString:  formData.Query,
		OperationName:  formData.Operation,
		VariableValues: formData.Variables,
	})

	return c.JSON(result)
}
