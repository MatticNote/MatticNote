package graphql

import (
	"context"
	"github.com/MatticNote/MatticNote/internal"
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

	currentUsr, ok := c.Locals(internal.LoginUserLocal).(*internal.LocalUserStruct)
	if ok {
		usrContext = context.WithValue(usrContext, "currentUser", currentUsr.Uuid)
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
