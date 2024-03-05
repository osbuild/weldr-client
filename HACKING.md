# Hacking on composer-cli

This new implementation of composer-cli is broken up into 2 major pieces. The
`./weldr/` package contains the library functions for communicating with the
API server. These functions are publicly accessible and should be documented so
that 3rd party users can use them to communicate with the osbuild-composer
server.

The `./cmd/composer-cli/` code is the user interface code, separated out into
commands and subcommands. They are for interactive use, and by scripts via the
`--json` argument that switches the output to JSON formated objects.


## Adding new library functions

The weldr library should contain functions, separated into descriptively named
files, like blueprints.go. The functions should handle communication with the
API server and return values to the caller using, if possible, regular Go data
types. eg. `ListBlueprints` returns a slice of strings.

Errors from lower functions should be returned to the called as errors, and API
error responses should be returned as `APIResponse` or slices of `APIErrorMsg`.
The type returned will depend on the response from the server.

The blueprints.go has examples of most of the types of responses you will
encounter, use it as a template for adding new functions.

Not every response can use a Go type. For those that need to return more
details use structs that match the API structs used in the server. We do not
want to import the actual server modules (and cannot because they are all
implemented as internal modules) so the structs should be placed into
`weldr/apischema.go`

Try to avoid using `interface{}` as much as possible. I have only used it for
returning the raw blueprint to the caller in the `GetBlueprints...` functions so
that it isn't tightly coupled to the actual blueprint format. The composer-cli
isn't manipulating the blueprint data, only passing it through, so it doesn't
need to know what's actually there. Blueprints are the thing most likely to
change, and they are not versioned like the API is so remaining flexible is
important.


## Unit tests for new functions

Functions, or parts of functions, that can be tested without an actual API
server should be written in the `weldr/unit_test.go` file, or if there are no
integration tests into a foo_test.go file where foo is the name of the file you
put the function in.


## Integration tests for library functions

Integration tests require an API server to be running, and should be placed
into a foo_test.go file alongside the code they are testing. At the top of the
file there should be a `// +build integration` before the package line, ans
separated from it with a blank line. See blueprints_test.go for an example.

These tests should be independent from each other. If there is server setup
that needs to be done it should be done in a way that works with a clean or a
dirty server. Usually it will be run against a clean server, as part of the
github action tests that are run when a PR is submitted. But don't depend on that.

If something needs to be setup you can do it in the test function itself, or in
the integration_test.go file if necessary. But try to keep changes there to a
minimum.


## Adding new cli commands

Mention Cobra, point to docs. Point to simple example like 'status show', tell how to hook it up the `rootCmd`
and how to write tests for it.

composer-cli uses the [Cobra](https://pkg.go.dev/github.com/spf13/cobra)
library to handle command parsing, options, and subcommands.

The `./cmd/composer-cli/` directory is structured so that top level commands go
into their own directory and package.  They import the root Cobra object from
`cmd/composer-cli/root` and attach themselves to that command using
`root.AddRootCommand()`. Then the subcommands attach themselves to this
command. See `cmd/composer-cli/blueprints/blueprints.go` for an example of a
subcommand.

Developer documentation should go into a doc.go file.

Commands are hooked into the root command parser by adding them to the include
list in `cmd/composer-cli/main.go` with a preceeding `_` so that the compiler
won't complain about unused symbols.


## Adding new subcommands

Subcommands should go into their own file, and be members of the command's
package. The file should only contain the functions needed to implement the
user interface for this command. Any helper functions or API related functions
that could be useful to other developers should go into the `weldr/` library.

## Error handling for cli commands

Commands should use the Cobra RunE method which returns an error to the root command handler.
If the command exits immediately on error it should return a
`root.ExecutionError` with a string describing the error, and embedding the `err` if there was one.
See `cmd/composer-cli/blueprints/list.go` for an example of this type of handling.

If there could be more than one error, eg. when processing a list of
blueprints, the command function should print them to `os.Stderr` as `ERROR: ...\n`
and return a blank `root.ExecutionError(cmd, "")` to the root handler. See
`cmd/composer-cli/blueprints/push.go` for an example.


## Unit tests for cli commands

The command line code under `./cmd/` is unit tested against the `./weldr/` library
code. They use a mock http client that is setup, by each test, to return the
expected JSON or TOML responses. See the tests in
`./cmd/composer-cli/blueprints/` for examples of how to set this up.

Basically, you setup a mock http client that stores the query it receives and
returns a canned response. This should not be too elaborate, otherwise you end
up mocking every detail of the server. The goal is to make sure the cmdline
command calls the right library function and that the response from the
function is parsed correctly.

The `weldr.MockClient` function that is used to handle the requests adds a `.Req`
element that can be examined in the tests. eg:

```
assert.Equal(t, "/api/v1/blueprints/list", mc.Req.URL.Path)
```

If the test returns an error the `MockClient` `DoFunc` needs to return the
request so that the `apiError` function can pass it to the callback function.
See `common_test.TestRequestMethods404` for an example.


# Running tests

## Running Unit Tests

Unit tests can be run by running `make test`.


## Running Integration Tests

Integration tests are run by a test binary. You can build the binary by running
`make integration` and then running it on a system with osbuild-composer running. Pass it -test.v to
output more verbose details about the tests being run.


