package handlers

import (
	"errors"
	"fmt"

	. "app/middlewares"
	"app/models"
	"app/store"
	"app/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserControlGet(c *fiber.Ctx) error {
	return c.Render("users", fiber.Map{
		CSRF:    c.Locals(CSRF),
		HTMX:    c.Locals(HTMX),
		"Users": models.GetAllUsers(),
	})
}

func UserEditGet(c *fiber.Ctx) error {
	if !c.Locals(HTMX).(bool) {
		return fiber.ErrNotFound
	}
	id := c.Params("id")
	u := new(models.User)
	err := models.GetUser(u, id)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return fiber.ErrNotFound
	}
	u.Password = ""
	return c.Render("useredit", fiber.Map{
		CSRF:        c.Locals(CSRF),
		models.USER: u,
	})
}

func UserEditPost(c *fiber.Ctx) error {
	if !c.Locals(HTMX).(bool) {
		return fiber.ErrNotFound
	}
	id := c.Params("id")
	f := new(models.UserForm)
	u := new(models.User)

	if err := c.BodyParser(f); err != nil {
		return err
	}

	err := models.GetUser(u, id)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return fiber.ErrNotFound
	}

	fmap := fiber.Map{
		CSRF:        c.Locals(CSRF),
		models.USER: u,
	}

	if message, ok := f.IsEmptyUsername(); ok {
		fmap["Message"] = message
		return c.Render("useredit", fmap)
	}
	if message, ok := f.IsPasswordsMatch(); ok {
		fmap["Message"] = message
		return c.Render("useredit", fmap)
	}
	if _, ok := f.IsPasswordsEmpty(); ok {
		f.Password, f.PasswordConfirm = u.Password, u.Password
	}
	if message, ok := f.IsPasswordsShort(); ok {
		fmap["Message"] = message
		return c.Render("useredit", fmap)
	}

	password, bcerr := utils.GetHash(f.Password)
	if bcerr != nil {
		return bcerr
	}

	u.Set(f.Username, string(password), f.IsAdmin())

	updateerr := models.UpdateUser(u)
	if updateerr.Error != nil {
		return err.Error
	}

	fmap["Success"] = fmt.Sprintf("Successful update user - %s", u.Username)
	return c.Render("useredit", fmap)
}

func UserCreateGet(c *fiber.Ctx) error {
	if !c.Locals(HTMX).(bool) {
		return fiber.ErrNotFound
	}
	return c.Render("usercreate", fiber.Map{CSRF: c.Locals(CSRF)})
}

func UserCreatePost(c *fiber.Ctx) error {
	f := new(models.UserForm)

	if err := c.BodyParser(f); err != nil {
		return err
	}
	u := models.User{Username: f.Username, Password: "", Admin: f.IsAdmin()}

	if message, ok := f.IsEmptyUsername(); ok {
		return c.Render("usercreate", fiber.Map{
			CSRF:      c.Locals(CSRF),
			"Message": message,
		})
	}
	if message, ok := f.CheckPasswordForCreate(); ok {
		return c.Render("usercreate", fiber.Map{
			CSRF:        c.Locals(CSRF),
			"Message":   message,
			models.USER: u,
		})
	}

	r := models.GetUserByUsername(&models.User{}, f.Username)
	if !errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return c.Render("usercreate", fiber.Map{
			CSRF:        c.Locals(CSRF),
			"Message":   fmt.Sprintf("User %s already exists", f.Username),
			models.USER: u,
		})
	}

	password, bcerr := utils.GetHash(f.Password)
	if bcerr != nil {
		return bcerr
	}
	u.Password = string(password)

	err := models.CreateUser(&u)

	if err.Error != nil {
		return err.Error
	}
	return c.Render("usercreate", fiber.Map{
		CSRF:        c.Locals(CSRF),
		"Success":   fmt.Sprintf("Succesful create user %s", u.Username),
		models.USER: u,
	})
}

func UserDeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	u := new(models.User)
	err := models.GetUser(u, id)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return fiber.ErrNotFound
	}
	deleterr := models.DeleteUser(u)
	if deleterr.Error != nil {
		return deleterr.Error
	}
	return c.Render("users", fiber.Map{
		CSRF:    c.Locals(CSRF),
		HTMX:    c.Locals(HTMX),
		"Users": models.GetAllUsers(),
	})
}

func LoginGet(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{CSRF: c.Locals(CSRF)})
}

func LoginPost(c *fiber.Ctx) error {
	f := new(models.User)
	u := new(models.User)

	if err := c.BodyParser(f); err != nil {
		return err
	}

	r := models.GetUserByUsername(u, f.Username)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return c.Render("login", fiber.Map{
			"Message": "Missmatch username or password",
			CSRF:      c.Locals(CSRF)},
		)
	}

	if err := utils.CompareHashAndPassword(u.Password, f.Password); err != nil {
		return c.Render("login", fiber.Map{
			"Message": "Missmatch username or password",
			CSRF:      c.Locals(CSRF)},
		)
	}

	session, err := store.Store.Get(c)
	if err != nil {
		return err
	}
	session.Set(store.AUTH_KEY, true)
	session.Set(store.USER_ID, u.ID)

	if sessionerr := session.Save(); sessionerr != nil {
		return err
	}
	return c.Redirect("/")
}

func LogoutGet(c *fiber.Ctx) error {
	session, err := store.Store.Get(c)
	if err != nil {
		return c.RedirectToRoute("login", fiber.Map{})
	}
	err = session.Destroy()
	if err != nil {
		return err
	}
	return c.RedirectToRoute("login", fiber.Map{})
}

func Health(c *fiber.Ctx) error {
	session, err := store.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("not auth")
	}
	auth := session.Get(store.AUTH_KEY)
	if auth != nil {
		return c.Status(fiber.StatusOK).SendString("auth")
	} else {
		return c.Status(fiber.StatusUnauthorized).SendString("not auth")
	}
}

func UserGet(c *fiber.Ctx) error {
	u := c.Locals(models.USER)
	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("login", fiber.Map{})
	}
	return c.Status(fiber.StatusOK).Render("index", fiber.Map{models.USER: user})
}
