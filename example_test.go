package rb_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	rb "github.com/gohandle/rb"
	"github.com/gohandle/rb/rbcore"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// Ctrl represents the C in MVC (model-view-conroller). It will have
// methods that create http.Handlers for different parts of your web application.
type Ctrl struct {

	// We embed our rb.App instance. It holds the configured framework.
	rb.App

	// It is common for the ctrl to embed other application dependencies. Such as
	// database connections, email senders, etc.
	db *sql.DB
}

// NewCtrl is the ctrl constructor. It takes all the dependencies and retruns the
// controller.
func NewCtrl(rba rb.App, db *sql.DB) *Ctrl {
	return &Ctrl{rba, db}
}

// HandleHome creates a http.Handler to handles requests for a Home page. It is can be
// called to create a Handler wherever that is needed: most commonly when building routes.
// The signature can be changed to any arguments, for example to allow configuration to be passed in.
func (c Ctrl) HandleHome(defaultGreeting string) http.Handler {

	// It is common to define a local type that holds all the state for handling requests
	// so it can be rendered.
	type Page struct {
		Greeting string
	}

	// The Action method is used to create the actual handler. It takes a function that will
	// be called on every request. It allows for returning an server errors and it takes a
	// rb.Ctx which provides common functionality
	return c.Action(func(c rb.Ctx) error {

		// This function will be called conncurrently. So make sure to declare any variables
		// locally to prevent race data races.
		p := Page{Greeting: defaultGreeting}

		// Render the template
		return c.Render(rb.Template("home.html", p))
	})
}

// HomeTmpl constains the jet template for rendering the home page. It is more common to
// load this from a filesystem. The render call in the action function determins what is
// passed in as the '.'. In this example
const HomeTmpl = `<h1>{{ .Greeting }}, world!</h1>`

func Example() {

	// rb uses popular packages in the Go ecosystem to provide the common functionality.
	// The developer is reponsible for initializing them and can be customized and configured
	// in any way desired.
	r := mux.NewRouter()
	tmpls := jet.NewInMemLoader()
	ss := sessions.NewCookieStore(make([]byte, 32))
	val := validator.New()
	fdec := form.NewDecoder()
	bdl := i18n.NewBundle(language.English)

	// for this example we'll set the template in memory
	tmpls.Set("/home.html", HomeTmpl)

	// The rb.Core abstracts the framework's external dependencies. The default core
	// builds on popular packages in the go ecosystem.
	core := rbcore.NewDefault(r, jet.NewSet(tmpls), fdec, val, ss, bdl)

	// we create the rb.App from the core, it also requires a zap Logger to
	// report any errors and debug info.
	logs, _ := zap.NewDevelopment()
	rba := rb.New(core, logs)

	// create te application dependencies, normally you would connect to a real database
	db := (*sql.DB)(nil)

	// create our application controller with the rb.App and our fake database
	ctrl := NewCtrl(rba, db)

	// configure a route with by calling the handler factory
	r.Handle("/", ctrl.HandleHome("hello"))

	// finally we can send a test request
	wr, req := httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(wr, req)

	fmt.Println(wr.Body.String())
	// Output: <h1>hello, world!</h1>

}
