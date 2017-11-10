package app

import "github.com/gin-gonic/gin"
import "net/http"
import v "gopkg.in/go-playground/validator.v9"
import "github.com/go-playground/universal-translator"
import "github.com/mcgaw/gusser/app/templates"

// Routes establies handlers for
// user http requests.
func Routes(e *gin.Engine) {

	e.GET("/register", func(c *gin.Context) {
		templates.Register(c.Writer, templates.RegisterForm{}, make(map[string]string))
	})

	e.POST("/register", func(c *gin.Context) {
		var form templates.RegisterForm
		c.ShouldBind(&form)

		_, err := createUser(form.Username, form.Password, form.DisplayName)
		if err != nil {
			if valErr, ok := err.(v.ValidationErrors); ok {
				t, _ := c.Get("translator")
				errors := valErr.Translate(t.(ut.Translator))
				templates.Register(c.Writer, form, errors)
			} else {
				c.AbortWithError(500, err)
			}
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	})
}

type User struct {
	Id          int    `db:"id"`
	Username    string `validate:"required,min=3"db:"username"`
	Password    string `validate:"required,min=8"db:"password"`
	DisplayName string `validate:"omitempty,min=3"db:"display_name"`
}

func createUser(username string, password string, displayName string) (*User, error) {

	user := &User{Username: username, Password: password, DisplayName: displayName}

	err := validator.Struct(user)
	if err != nil {
		if valErr, ok := err.(v.ValidationErrors); ok {
			return nil, valErr
		} else if err != nil {
			panic(err)
		}
	}

	db.MustExec("insert into user (username, password, display_name) values (?,?,?)",
		username, password, displayName)
	return user, nil
}

func getUser(username string) *User {
	user := User{}
	err := db.Get(&user, "select * from user where username=$1", username)
	if err != nil {
		panic(err)
	}
	return &user
}
