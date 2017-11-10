package app

import (
	"errors"

	// Import to link the C sqlite driver.
	_ "github.com/mattn/go-sqlite3"

	en_translations "gopkg.in/go-playground/validator.v9/translations/en"

	"github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/mcgaw/gusser/app/templates"

	"github.com/jmoiron/sqlx"

	v "gopkg.in/go-playground/validator.v9"
)

var db *sqlx.DB
var validator *v.Validate

func init() {
	validator = v.New()
	dbase, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	Create(dbase)
	db = dbase
}

// Validation is a middleware function to make a suitable
// translator available on the Gin context (currently just
// supports one Translator).
func Validation(t ut.Translator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("translator", t)
	}
}

/*
func addLocalMessages(uT *ut.UniversalTranslator) {
	transErr := uT.Import(ut.FormatJSON, "app/translations")
	if transErr != nil {
		panic(transErr)
	}
	transErr = uT.VerifyTranslations()
	if transErr != nil {
		panic(transErr)
	}
}
*/

func Start() {
	var r *gin.Engine
	r = Setup("")
	r.Run()
}

func Setup(mode string) *gin.Engine {
	r := gin.Default()

	if mode != "" {
		gin.SetMode(mode)
	}

	// Setup English translation, including support for validaion keys.
	enLoc := en.New()
	uT := ut.New(enLoc, enLoc)
	enTrans, _ := uT.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validator, enTrans)
	r.Use(Validation(enTrans))

	// Setup up route handlers.
	Routes(r)

	r.GET("/", func(c *gin.Context) {
		templates.Index(c.Writer)
	})
	return r
}

// Params is a template function that can be called in a
// pipe to simulate parameter passing (to different templates)
func Params(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
