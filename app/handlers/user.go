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

var (
	PASS_LEN int    = 6
	User     string = "User"
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
	u.Password = ""
	c.Locals(User, u)

	return c.Next()
}

func SuperUserMiddleware(c *fiber.Ctx) error {
	u := c.Locals(User)
	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("auth-login", fiber.Map{})
	}
	if !user.Admin {
		return fiber.ErrNotFound
	}
	return c.Next()
}

func LoginPageGet(c *fiber.Ctx) error {
	return c.Render("templates/login", fiber.Map{"csrf": c.Locals("csrf")})
}

func LoginPost(c *fiber.Ctx) error {
	f := new(models.User)
	u := new(models.User)

	if err := c.BodyParser(f); err != nil {
		return err
	}

	r := models.GetUserByUsername(u, f.Username)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		// if c.Locals(Htmx).(bool) {
		// 	return c.Status(fiber.StatusOK).SendString("<div>Missmatch username or password</biv>")
		// }
		return c.Status(fiber.StatusOK).Render("templates/login", fiber.Map{
			"Message": "Missmatch username or password",
			Csrf:      c.Locals(Csrf)},
		)
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(f.Password)); err != nil {
	if err := utils.CompareHashAndPassword(u.Password, f.Password); err != nil {
		return c.Status(fiber.StatusOK).Render("templates/login", fiber.Map{
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

func CreateNewUserGet(c *fiber.Ctx) error {
	return c.Render("templates/userform", fiber.Map{Csrf: c.Locals(Csrf)})
}

func CreateNewUserPost(c *fiber.Ctx) error {
	f := new(models.UserForm)

	if err := c.BodyParser(f); err != nil {
		return err
	}

	if f.Username == "" {
		return c.Status(fiber.StatusBadRequest).Render("templates/userform", fiber.Map{"Message": "Username is required."})
	}

	if f.Password == "" || f.PasswordConfirm == "" {
		return c.Status(fiber.StatusBadRequest).Render("templates/userform", fiber.Map{"Message": "Password is required."})
	}

	if f.Password != f.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).Render("templates/userform", fiber.Map{"Message": "The passwords don't match.", "Username": f.Username})
	}

	if len([]rune(f.Password)) < PASS_LEN {
		return c.Status(fiber.StatusBadRequest).Render("templates/userform", fiber.Map{"Message": fmt.Sprintf("The minimum len of password is %d", PASS_LEN), "Username": f.Username})
	}

	r := models.GetUserByUsername(&models.User{}, f.Username)
	if !errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusBadRequest).Render("templates/userform", fiber.Map{"Message": fmt.Sprintf("User %s already exists", f.Username)})
	}

	password, bcerr := utils.GetHash(f.Password)
	if bcerr != nil {
		return bcerr
	}

	u := models.User{Username: f.Username, Password: string(password), Admin: false}

	if f.Admin != "" {
		u.Admin = true
	}

	err := models.CreateUser(&u)

	if err.Error != nil {
		return err.Error
	}
	return c.Status(fiber.StatusOK).Render("templates/userform", fiber.Map{"Message": fmt.Sprintf("Succesful create user %s", u.Username)})
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
	u := c.Locals(User)
	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("auth-login", fiber.Map{})
	}
	return c.Status(fiber.StatusOK).Render("templates/index", fiber.Map{User: user})
}
