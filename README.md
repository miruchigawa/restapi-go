# README

OpenSource RestAPI with downloader and anime features.

## Getting started

Make sure that you're in the root of the project directory, fetch the dependencies with `go mod tidy`, then run the application using `go run ./cmd/api`:

```
$ go mod tidy
$ go run ./cmd/api
```

If you make a request to the `GET /status` endpoint using `curl` you should get a response like this:

```
$ curl -i localhost:4444/status
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 09 May 2022 20:46:37 GMT
Content-Length: 23

{
    "Status": "OK",
}
```

You can also start the application with live reload support by using the `run` task in the `Makefile`:

```
$ make run
```

## Project structure

Everything in the codebase is designed to be editable. Feel free to change and adapt it to meet your needs.

|     |     |
| --- | --- |
| **`assets`** | Contains the non-code assets for the application. |
| `↳ assets/emails/` | Contains email templates. |
| `↳ assets/migrations/` | Contains SQL migrations. |
| `↳ assets/efs.go` | Declares an embedded filesystem containing all the assets. |

|     |     |
| --- | --- |
| **`cmd/api`** | Your application-specific code (handlers, routing, middleware, helpers) for dealing with HTTP requests and responses. |
| `↳ cmd/api/errors.go` | Contains helpers for managing and responding to error conditions. |
| `↳ cmd/api/handlers.go` | Contains your application HTTP handlers. |
| `↳ cmd/api/helpers.go` | Contains helper functions for common tasks. |
| `↳ cmd/api/main.go` | The entry point for the application. Responsible for parsing configuration settings initializing dependencies and running the server. Start here when you're looking through the code. |
| `↳ cmd/api/middleware.go` | Contains your application middleware. |
| `↳ cmd/api/routes.go` | Contains your application route mappings. |
| `↳ cmd/api/server.go` | Contains a helper functions for starting and gracefully shutting down the server. |

|     |     |
| --- | --- |
| **`internal`** | Contains various helper packages used by the application. |
| `↳ internal/database/` | Contains your database-related code (setup, connection and queries). |
| `↳ internal/env` | Contains helper functions for reading configuration settings from environment variables. |
| `↳ internal/funcs/` | Contains custom template functions. |
| `↳ internal/request/` | Contains helper functions for decoding JSON requests. |
| `↳ internal/response/` | Contains helper functions for sending JSON responses. |
| `↳ internal/smtp/` | Contains a SMTP sender implementation. |
| `↳ internal/validator/` | Contains validation helpers. |
| `↳ internal/version/` | Contains the application version number definition. |

## Configuration settings

Configuration settings are managed via environment variables, with the environment variables read into your application in the `run()` function in the `main.go` file.

You can try this out by setting a `HTTP_PORT` environment variable to configure the network port that the server is listening on:

```
$ export HTTP_PORT="9999"
$ go run ./cmd/api
```

Feel free to adapt the `run()` function to parse additional environment variables and store their values in the `config` struct. The application uses helper functions in the `internal/env` package to parse environment variable values or return a default value if no matching environment variable is set. It includes `env.GetString()`, `env.GetInt()` and `env.GetBool()` functions for reading string, integer and bool values from environment variables. Again, you can add any additional helper functions that you need.

## Creating new handlers

Handlers are defined as `http.HandlerFunc` methods on the `application` struct. They take the pattern:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    // Your handler logic...
}
```

Handlers are defined in the `cmd/api/handlers.go` file. For small applications, it's fine for all handlers to live in this file. For larger applications (10+ handlers) you may wish to break them out into separate files.

## Handler dependencies

Any dependencies that your handlers have should be initialized in the `run()` function `cmd/api/main.go` and added to the `application` struct. All of your handlers, helpers and middleware that are defined as methods on `application` will then have access to them.

You can see an example of this in the `cmd/api/main.go` file where we initialize a new `logger` instance and add it to the `application` struct.

## Creating new routes

[Gorilla mux](https://github.com/gorilla/mux) is used for routing. Routes are defined in the `routes()` method in the `cmd/api/routes.go` file. For example:

```
func (app *application) routes() http.Handler {
    mux := mux.NewRouter()

    mux.HandleFunc("/your/path", app.yourHandler).Methods("GET")

    mux.Get("/your/path", app.yourHandler)

    return mux
}
```

For more information about Gorilla mux and example usage, please see the [official documentation](https://github.com/gorilla/mux).

## Adding middleware

Middleware is defined as methods on the `application` struct in the `cmd/api/middleware.go` file. Feel free to add your own. They take the pattern:

```
func (app *application) yourMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Your middleware logic...
        next.ServeHTTP(w, r)
    })
}
```

## Sending JSON responses

JSON responses and a specific HTTP status code can be sent using the `response.JSON()` function. The `data` parameter can be any JSON-marshalable type.

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"hello":  "world"}

    err := response.JSON(w, http.StatusOK, data)
    if err != nil {
        app.serverError(w, r, err)
    }
}
```

Specific HTTP headers can optionally be sent with the response too:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"hello":  "world"}

    headers := make(http.Header)
    headers.Set("X-Server", "Go")

    err := response.JSONWithHeaders(w, http.StatusOK, data, headers)
    if err != nil {
        app.serverError(w, r, err)
    }
}
```

## Parsing JSON requests

HTTP requests containing a JSON body can be decoded using the `request.DecodeJSON()` function. For example, to decode JSON into an `input` struct:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name string `json:"Name"`
        Age  int    `json:"Age"`
    }

    err := request.DecodeJSON(w, r, &input)
    if err != nil {
        app.badRequest(w, r, err)
        return
    }

    ...
}
```

Note: The target decode destination passed to `request.DecodeJSON()` (which in the example above is `&input`) must be a non-nil pointer.

The `request.DecodeJSON()` function returns friendly, well-formed, error messages that are suitable to be sent directly to the client using the `app.badRequest()` helper.

There is also a `request.DecodeJSONStrict()` function, which works in the same way as `request.DecodeJSON()` except it will return an error if the request contains any JSON fields that do not match a name in the the target decode destination.

## Validating JSON requests

The `internal/validator` package includes a simple (but powerful) `validator.Validator` type that you can use to carry out validation checks.

Extending the example above:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name      string              `json:"Name"`
        Age       int                 `json:"Age"`
        Validator validator.Validator `json:"-"`
    }

    err := request.DecodeJSON(w, r, &input)
    if err != nil {
        app.badRequest(w, r, err)
        return
    }

    input.Validator.CheckField(input.Name != "", "Name", "Name is required")
    input.Validator.CheckField(input.Age != 0, "Age", "Age is required")
    input.Validator.CheckField(input.Age >= 21, "Age", "Age must be 21 or over")

    if input.Validator.HasErrors() {
        app.failedValidation(w, r, input.Validator)
        return
    }

    ...
}
```

The `app.failedValidation()` helper will send a `422` status code along with any validation error messages. For the example above, the JSON response will look like this:

```
{
    "FieldErrors": {
        "Age": "Age must be 21 or over",
        "Name": "Name is required"
    }
}
```

In the example above we use the `CheckField()` method to carry out validation checks for specific fields. You can also use the `Check()` method to carry out a validation check that is _not related to a specific field_. For example:

```
input.Validator.Check(input.Password == input.ConfirmPassword, "Passwords do not match")
```

The `validator.AddError()` and `validator.AddFieldError()` methods also let you add validation errors directly:

```
input.Validator.AddFieldError("Email", "This email address is already taken")
input.Validator.AddError("Passwords do not match")
```

The `internal/validator/helpers.go` file also contains some helper functions to simplify validations that are not simple comparison operations.

|     |     |
| --- | --- |
| `NotBlank(value string)` | Check that the value contains at least one non-whitespace character. |
| `MinRunes(value string, n int)` | Check that the value contains at least n runes. |
| `MaxRunes(value string, n int)` | Check that the value contains no more than n runes. |
| `Between(value, min, max T)` | Check that the value is between the min and max values inclusive. |
| `Matches(value string, rx *regexp.Regexp)` | Check that the value matches a specific regular expression. |
| `In(value T, safelist ...T)` | Check that a value is in a 'safelist' of specific values. |
| `AllIn(values []T, safelist ...T)` | Check that all values in a slice are in a 'safelist' of specific values. |
| `NotIn(value T, blocklist ...T)` | Check that the value is not in a 'blocklist' of specific values. |
| `NoDuplicates(values []T)` | Check that a slice does not contain any duplicate (repeated) values. |
| `IsEmail(value string)` | Check that the value has the formatting of a valid email address. |
| `IsURL(value string)` | Check that the value has the formatting of a valid URL. |

For example, to use the `Between` check your code would look similar to this:

```
input.Validator.CheckField(validator.Between(input.Age, 18, 30), "Age", "Age must between 18 and 30")
```

Feel free to add your own helper functions to the `internal/validator/helpers.go` file as necessary for your application.

## Working with the database

This codebase is set up to use SQLite3 with the [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) driver. The data is stored in a `db.sqlite` file in the project root, but you can change this by setting a different DSN (datasource name) in the `DB_DSN` environment variable , or by adapting the default value in `run()`.

The codebase is also configured to use [jmoiron/sqlx](https://github.com/jmoiron/sqlx), so you have access to the whole range of sqlx extensions as well as the standard library `Exec()`, `Query()` and `QueryRow()` methods .

The database is available to your handlers, middleware and helpers via the `application` struct. If you want, you can access the database and carry out queries directly. For example:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    _, err := app.db.Exec("INSERT INTO people (name, age) VALUES ($1, $2)", "Alice", 28)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    ...
}
```

Generally though, it's recommended to isolate your database logic in the `internal/database` package and extend the `DB` type to include your own methods. For example, you could create a `internal/database/people.go` file containing code like:

```
type Person struct {
    ID    int    `db:"id"`
    Name  string `db:"name"`
    Age   int    `db:"age"`
}

func (db *DB) NewPerson(name string, age int) error {
    _, err := db.Exec("INSERT INTO people (name, age) VALUES ($1, $2)", name, age)
    return err
}

func (db *DB) GetPerson(id int) (Person, error) {
    var person Person
    err := db.Get(&person, "SELECT * FROM people WHERE id = $1", id)
    return person, err
}
```

And then call this from your handlers:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    err := app.db.NewPerson("Alice", 28)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    ...
}
```

## Managing SQL migrations

The `Makefile` in the project root contains commands to easily create and work with database migrations:

|     |     |
| --- | --- |
| `$ make migrations/new name=add_example_table` | Create a new database migration in the `assets/migrations` folder. |
| `$ make migrations/up` | Apply all up migrations. |
| `$ make migrations/down` | Apply all down migrations. |
| `$ make migrations/goto version=N` | Migrate up or down to a specific migration (where N is the migration version number). |
| `$ make migrations/force version=N` | Force the database to be specific version without running any migrations. |
| `$ make migrations/version` | Display the currently in-use migration version. |

Hint: You can run `$ make help` at any time for a reminder of these commands.

These `Makefile` tasks are simply wrappers around calls to the `github.com/golang-migrate/migrate/v4/cmd/migrate` tool. For more information, please see the [official documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

By default all 'up' migrations are automatically run on application startup using embeded files from the `assets/migrations` directory. You can disable this by setting the `DB_AUTOMIGRATE` environment variable to `false`.

Important: If you change your database location from the default `db.sqlite` file then you will need to also update the `Makefile` tasks to reference the new location.

## Logging

Leveled logging is supported using the [slog](https://pkg.go.dev/log/slog) and [tint](https://github.com/lmittmann/tint) packages.

By default, a logger is initialized in the `main()` function. This logger writes all log messages above `Debug` level to `os.Stdout`.

```
logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
```

Feel free to customize this further as necessary.

Also note: Any messages that are automatically logged by the Go `http.Server` are output at the `Warn` level.

## Sending emails

The application is configured to support sending of emails via SMTP.

Email templates should be defined as files in the `assets/emails` folder. Each file should contain named templates for the email subject, plaintext body and — optionally — HTML body.

```
{{define "subject"}}Example subject{{end}}

{{define "plainBody"}}
This is an example body
{{end}}

{{define "htmlBody"}}
<!doctype html>
<html>
    <head>
        <meta name="viewport" content="width=device-width" />
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    </head>
    <body>
        <p>This is an example body</p>
    </body>
</html>
{{end}}
```

A further example can be found in the `assets/emails/example.tmpl` file. Note that your email templates automatically have access to the custom template functions defined in the `internal/funcs` package.

Emails can be sent from your handlers using `app.mailer.Send()`. For example, to send an email to `alice@example.com` containing the contents of the `assets/emails/example.tmpl` file:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    data := map[string]any{"Name": "Alice"}

    err := app.mailer.Send("alice@example.com", data, "example.tmpl")
    if err != nil {
        app.serverError(w, r, err)
        return
    }

   ...
}
```

Note: The second parameter to `Send()` should be a map or struct containing any dynamic data that you want to render in the email template.

The SMTP host, port, username, password and sender details can be configured using the `SMTP_HOST` environment variable, `SMTP_PORT` environment variable, `SMTP_USERNAME` environment variable, `SMTP_PASSWORD` environment variable, and `SMTP_FROM` environment variable or by adapting the default values in `cmd/api/main.go`.

You may wish to use [Mailtrap](https://mailtrap.io/) or a similar tool for development purposes.

## Error notifications

The application supports sending alerts for runtime errors to an admin email address. You can enable this by setting the `NOTIFICATIONS_EMAIL` environment variable to a valid email address.

If you don't set a value for this, then no notifications will be sent (but the errors will still be logged).

Notifications will only be sent for any errors that are encountered as part of a request-response cycle (i.e. whenever a user sees an '500 Internal Server Error' response). Notifications are not sent for any errors that occur when starting or shutting down the application, so it's important to still use an uptime monitoring service in production.

The code for this functionality is in the `sendErrorNotification()` method (in the `cmd/api/errors.go` file) and the email template for the notification is located at `assets/emails/error-notification.tmpl`.

You'll need to make sure that the application is [configured using valid SMTP credentials](#sending-emails) in order to send the email.

## Custom template functions

Custom template functions are defined in `internal/funcs/funcs.go` and are automatically made available to your

email templates when you use `app.mailer.Send()`
.

The following custom template functions are already included by default:

|     |     |
| --- | --- |
| `now` | Returns the current time. |
| `timeSince arg1` | Returns the time elapsed since arg1. |
| `timeUntil arg2` | Returns the time until arg1. |
| `formatTime arg1 arg2` | Returns the time arg2 as formatted using the pattern arg1. |
| `approxDuration arg1` | Returns the approximate duration of arg1 in a 'human-friendly' format ("3 seconds", "2 months", "5 years") etc. |
| `uppercase arg1` | Returns arg1 converted to uppercase. |
| `lowercase arg1` | Returns arg1 converted to lowercase. |
| `pluralize arg1 arg2 arg3` | If arg1 equals 1 then return arg2, otherwise return arg3. |
| `slugify arg1` | Returns the lowercase of arg1 with all non-ASCII characters and punctuation removed (expect underscores and hyphens). Whitespaces are also replaced with a hyphen. |
| `safeHTML arg1` | Output the verbatim value of arg1 without escaping the content. This should only be used when arg1 is from a trusted source. |
| `join arg1 arg2` | Returns the values in slice arg1 joined using the separator arg2. |
| `incr arg1` | Increments arg1 by 1. |
| `decr arg1` | Decrements arg1 by 1. |
| `formatInt arg1` | Returns arg1 formatted with commas as the thousands separator. |
| `formatFloat arg1 arg2` | Returns arg1 rounded to arg2 decimal places and formatted with commas as the thousands separator. |
| `yesno arg1` | Returns "Yes" if arg1 is true, or "No" if arg1 is false. |
| `urlSetParam arg1 arg2 arg3` | Returns the URL arg1 with the key arg2 and value arg3 added to the query string parameters. |
| `urlDelParam arg1 arg2` | Returns the URL arg1 with the key arg2 (and corresponding value) removed from the query string parameters. |

To add another custom template function, define the function in `internal/funcs/funcs.go` and add it to the `TemplateFuncs` map. For example:

```
var TemplateFuncs = template.FuncMap{
    ...
    "yourFunction": yourFunction,
}

func yourFunction(s string) (string, error) {
    // Do something...
}
```

## Admin tasks

The `Makefile` in the project root contains commands to easily run common admin tasks:

|     |     |
| --- | --- |
| `$ make tidy` | Format all code using `go fmt` and tidy the `go.mod` file. |
| `$ make audit` | Run `go vet`, `staticheck`, `govulncheck`, execute all tests and verify required modules. |
| `$ make test` | Run all tests. |
| `$ make test/cover` | Run all tests and outputs a coverage report in HTML format. |
| `$ make build` | Build a binary for the `cmd/api` application and store it in the `/tmp/bin` folder. |
| `$ make run` | Build and then run a binary for the `cmd/api` application. |
| `$ make run/live` | Build and then run a binary for the `cmd/api` application (uses live reloading). |

## Live reload

When you use `make run/live` to run the application, the application will automatically be rebuilt and restarted whenever you make changes to any files with the following extensions:

```
.go
.tpl, .tmpl, .html
.css, .js, .sql
.jpeg, .jpg, .gif, .png, .bmp, .svg, .webp, .ico
```

Behind the scenes the live reload functionality uses the [cosmtrek/air](https://github.com/cosmtrek/air) tool. You can configure how it works (including which file extensions and folders are watched for changes) by editing the `Makefile` file.

## Running background tasks

A `backgroundTask()` helper is included in the `cmd/api/helpers.go` file. You can call this in your handlers, helpers and middleware to run any logic in a separate background goroutine. This useful for things like sending emails, or completing slow-running jobs.

You can call it like so:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    app.backgroundTask(r, func() error {
        // The logic you want to execute in a background task goes here.
        // It should return an error, or nil.
        err := doSomething()
        if err != nil {
            return err
        }

        return nil
    })

    ...
}
```

Using the `backgroundTask()` helper will automatically recover any panics in the background task logic, and when performing a graceful shutdown the application will wait for any background tasks to finish running before it exits.

## Application version

The application version number is generated automatically based on your latest version control system revision number. If you are using Git, this will be your latest Git commit hash. It can be retrieved by calling the `version.Get()` function from the `internal/version` package.

Important: The version control system revision number will only be available after you have initialized version control for your repository (e.g. run `$ git init`) AND application is built using `go build`. If you run the application using `go run` then `version.Get()` will return the string `"unavailable"`.

## Changing the module path

The module path is currently set to `miruchigawa.moe/restapi`. If you want to change this please find and replace all instances of `miruchigawa.moe/restapi` in the codebase with your own module path.
