package handlers

import (
	"errors"
	"fmt"

	"auth/app/models"
	"auth/app/store"
	"auth/app/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthMiddleware(c *fiber.Ctx) error {
	session, err := store.Store.Get(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("auth-login", fiber.Map{})
	}

	if session.Get(store.AUTH_KEY) == nil {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("auth-login", fiber.Map{})
	}

	uid := session.Get(store.USER_ID)
	if uid == nil {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("auth-login", fiber.Map{})
	}

	u := new(models.User)
	if models.GetUser(u, fmt.Sprint(uid)); err != nil {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("auth-login", fiber.Map{})
	}

	c.Locals(models.USER, u)

	return c.Next()
}

func SuperUserMiddleware(c *fiber.Ctx) error {
	u := c.Locals(models.USER)
	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("auth-login", fiber.Map{})
	}
	if !user.Admin {
		return fiber.ErrNotFound
	}
	return c.Next()
}

func UserControlGet(c *fiber.Ctx) error {
	users := models.GetAllUsers()
	return c.Render("users", fiber.Map{
		Csrf:    c.Locals(Csrf),
		Htmx:    c.Locals(Htmx),
		"Users": users,
	})
}

func UserEditGet(c *fiber.Ctx) error {
	if !c.Locals(Htmx).(bool) {
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
		Csrf:   c.Locals(Csrf),
		"User": u,
	})
}

func UserEditPost(c *fiber.Ctx) error {
	if !c.Locals(Htmx).(bool) {
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
		Csrf:        c.Locals(Csrf),
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
		f.Password = u.Password
		f.PasswordConfirm = u.Password
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

func CreateNewUserGet(c *fiber.Ctx) error {
	if !c.Locals(Htmx).(bool) {
		return fiber.ErrNotFound
	}
	return c.Render("usercreate", fiber.Map{Csrf: c.Locals(Csrf)})
}

func CreateNewUserPost(c *fiber.Ctx) error {
	f := new(models.UserForm)

	if err := c.BodyParser(f); err != nil {
		return err
	}

	if message, ok := f.IsEmptyUsername(); ok {
		return c.Render("usercreate", fiber.Map{
			Csrf:      c.Locals(Csrf),
			"Message": message,
		})
	}
	if message, ok := f.CheckPasswordForCreate(); ok {
		return c.Render("usercreate", fiber.Map{
			Csrf:       c.Locals(Csrf),
			"Message":  message,
			"Username": f.Username})
	}

	r := models.GetUserByUsername(&models.User{}, f.Username)
	if !errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return c.Render("usercreate", fiber.Map{
			Csrf:      c.Locals(Csrf),
			"Message": fmt.Sprintf("User %s already exists", f.Username),
		})
	}

	password, bcerr := utils.GetHash(f.Password)
	if bcerr != nil {
		return bcerr
	}

	u := models.User{Username: f.Username, Password: string(password), Admin: f.IsAdmin()}

	err := models.CreateUser(&u)

	if err.Error != nil {
		return err.Error
	}
	return c.Render("usercreate", fiber.Map{
		Csrf:      c.Locals(Csrf),
		"Success": fmt.Sprintf("Succesful create user %s", u.Username),
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
	users := models.GetAllUsers()
	return c.Render("users", fiber.Map{
		Csrf:    c.Locals(Csrf),
		Htmx:    c.Locals(Htmx),
		"Users": users,
	})
}

func LoginPageGet(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{Csrf: c.Locals(Csrf)})
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
			Csrf:      c.Locals(Csrf)},
		)
	}

	if err := utils.CompareHashAndPassword(u.Password, f.Password); err != nil {
		return c.Render("login", fiber.Map{
			"Message": "Missmatch username or password",
			Csrf:      c.Locals(Csrf)},
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
		return c.RedirectToRoute("auth-login", fiber.Map{})
	}
	err = session.Destroy()
	if err != nil {
		return err
	}
	return c.RedirectToRoute("auth-login", fiber.Map{})
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
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("auth-login", fiber.Map{})
	}
	return c.Status(fiber.StatusOK).Render("templates/index", fiber.Map{models.USER: user})
}

// [HTMX Middleware]
// if c.Locals(Htmx).(bool) {
// 	return c.Status(fiber.StatusOK).SendString("<div>Missmatch username or password</biv>")
// }
