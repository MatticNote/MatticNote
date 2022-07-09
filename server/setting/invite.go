package setting

import (
	"errors"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/MatticNote/MatticNote/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	"time"
)

type inviteCreateFormStruct struct {
	Count        int `form:"count"`
	ExpiredAfter int `form:"expired_after"`
}

type inviteDeleteFormStruct struct {
	ID string `form:"id"`
}

func settingInviteGet(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)

	if config.Config.System.RegistrationMode != 1 || (config.Config.System.InvitePermission == 0 && !user.IsAdmin) || (config.Config.System.InvitePermission == 1 && !(user.IsAdmin || user.IsModerator)) {
		return fiber.ErrForbidden
	}

	inviteList, err := account.GetUserInviteList(user.ID)
	if err != nil {
		return err
	}

	var newInvite *schemas.UserInvite
	newInvite, ok := c.Locals("createdInvite").(*schemas.UserInvite)
	if !ok {
		newInvite = nil
	}

	deletedInvite, ok := c.Locals("deletedInvite").(bool)
	if !ok {
		deletedInvite = false
	}

	return c.Render(
		"setting/invite",
		fiber.Map{
			"inviteList":    inviteList,
			"csrfName":      common.CSRFFormName,
			"csrfToken":     c.Locals(common.CSRFContextKey),
			"newInvite":     newInvite,
			"deletedInvite": deletedInvite,
		},
		"setting/_layout",
	)
}

func settingInvitePost(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)

	if config.Config.System.RegistrationMode != 1 || (config.Config.System.InvitePermission == 0 && !user.IsAdmin) || (config.Config.System.InvitePermission == 1 && !(user.IsAdmin || user.IsModerator)) {
		return fiber.ErrForbidden
	}

	form := new(actionFormStruct)
	if err := c.BodyParser(form); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	switch form.Action {
	case "create_invite":
		formCreate := new(inviteCreateFormStruct)
		if err := c.BodyParser(formCreate); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		var expiredAt *time.Time = nil
		if formCreate.ExpiredAfter > 0 {
			expiredAtRaw := time.Now().Add(time.Duration(formCreate.ExpiredAfter) * time.Second)
			expiredAt = &expiredAtRaw
		}

		newInvite, err := account.CreateInvite(&user.ID, uint(formCreate.Count), expiredAt)
		if err != nil {
			return err
		}

		c.Locals("createdInvite", newInvite)
		return settingInviteGet(c)
	case "delete_invite":
		formDelete := new(inviteDeleteFormStruct)

		if err := c.BodyParser(formDelete); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		id, err := ksuid.Parse(formDelete.ID)
		if err != nil {
			return fiber.ErrBadRequest
		}

		invite, err := account.GetInvite(id)
		if err != nil {
			if errors.Is(err, account.ErrInviteCodeNotFound) {
				return fiber.ErrBadRequest
			} else {
				return nil
			}
		}

		if invite.Owner == nil || (invite.Owner != nil && user.ID != *invite.Owner) {
			return fiber.ErrForbidden
		}

		err = account.DeleteInvite(id)
		if err != nil {
			return err
		}

		c.Locals("deletedInvite", true)
		return settingInviteGet(c)
	default:
		return fiber.ErrBadRequest
	}
}
