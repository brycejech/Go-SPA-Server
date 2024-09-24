# Go SPA Server

Go HTTP server for deploying modern single page applications (SPAs) as a single executable with configurable runtime
environment variables.

This application is built completely on the Go standard library and has no dependencies of any kind.

---

## Important notes

### File caching

This server pre-processes all of the embedded SPA files (included static assets) and caches them in memory. This gives a
performance boost in terms of response times (we only talk to the embedded filesystem at start time, swap bundled env
variables with runtime env variables, and pre-compute http response headers) but comes at the cost of increased memory
utilization.

This caching should cause memory utilization once started to be approximately double the size of the bundled SPA files.
If your SPA is aggressively large and your system resources are limited, this server may not be the best choice for you.

### Command line arguments

Because we chose not to use any CLI helper libraries like [Cobra](https://github.com/spf13/cobra), CLI arguments must be
provided precisely as documented. Arguments that are not provided will print warnings alongside default values.

### HTTP content types

In the Go standard library, there does not exist a great mechanism to provide robust, reliable, and trustworthy
content-type detection.

The `net/http` package provides a `DetectContentType()` function that examines file contents to perform content-type
detection but it does not support a fairly large number of common file types.

The `mime` package provides a `TypeByExtension(ext string)` function that returns content types based on file
extensions. It supports a greater number of content types (to varying degrees depending on execution OS) but,
importantly, _does not examine file contents_ to determine the content-type of the file.

Because the `mime` package provides a higher number of content types, we have decided to use it over the `net/http`
package. This does, however, mean that you, the user of this application, must have a higher degree of trust in the
files that you are serving in that their contents are what they report to be based on their extensions.

**Important:**

All files for which a content-type cannot be determined by extension will default to `application/octet-stream`.
Browsers will treat this as a binary file, and, if not initiated by an async mechanism such as `fetch`, download the
file for the user.

### All files in the artifact directory are embedded

The `go:embed` compiler directive is the mechanism that allows the SPA files to be bundled into the final executable. By
default, `go:embed` ignores files that are commonly "hidden" on some operating systems such as files beginning with `.`
or `_`, e.g. `.env` or `_hidden.txt`.

In order to support the swapping of stand-in, compiled build-time variables with runtime variables and because most
applications will utilize a `.env` file of some sort, we use the supplemental `:all` directive to instruct `go:embed` to
embed _all_ files within the bundled SPA.

This comes with the implications that **A)** you must include your build-time environment file within your SPA build
output and **B)** that file will be discoverable/downloadable by end users of the application. The latter of these two
implications should not be of concern, however, because the values for the build-time variables should be random v4
UUIDs and not represent anything meaningful. Regardless, this server intends to ship client-side code that should never
have any important secrets anyway.

---

## Building the server

Build your SPA with the build tool of your choosing into an output directory named `artifact/`.

Values for your environment _at build time_ should all be unique values, preferably v4 UUIDs. The build time environment
file _must be present in the output `artifact/` directory_. This file may need to be manually moved into the `artifact/`
directory after the fact as many SPA compilers/bundlers will not move it there by default.

Once built, move the `artifact/` directory into the top level directory for this application.

To build the final binary _for the current OS and CPU architecture_, use the following `go build` command:
`CGO_ENABLED=0 go build -a -o server -ldflags="-w -s" .`.

If you need to compile for a different OS/arch target, consult the
[Go docs on environment variables](https://pkg.go.dev/cmd/go#hdr-Environment_variables) for the `GOOS` and `GOARCH`
environment variables you will need.

Running this build command will output a single executable named `./server` that includes the bundled SPA files.

---

## Running the server

Once compiled, the server can be deployed and run by setting the runtime environment variables for the SPA and providing
a few command line arguments.

**Example:**

```sh
ENV_1='my env value' ENV_2='another env value' ./server --port=8000 --embeddedEnvFile='.env.dist' --requireAllVars=true
```

---

## Server arguments

The final build requires a few command line arguments to start: `port`, `staticDir`, `embeddedEnvFile`,
`requireAllVars`.

### port

`port` specifies the port that the HTTP server listens on

`--port=8000`

### staticDir

`staticDir` names the sub-directory of the embedded `artifact/` directory in which static assets are stored. If a
`staticDir` is provided, any requests for files within it that are not found will give an HTTP 404 response.

For all requests outside of the `staticDir`, any files that are not found are served the `index.html` file content.

`--staticDir=static`

### embeddedEnvFile

`embeddedEnvFile` indicates the name of the environment file that was used in the initial bundling and is included in
the `artifact/` directory. If an `embeddedEnvFile` is provided but not found in the embedded SPA files, the server will
error out and refuse to start.

`--embeddedEnvFile='.env.dist'`

### requireAllVars

`requireAllVars` is a boolean value that, when true, will verify that all variables present in the `embeddedEnvFile` are
also present in the runtime environment. If set to true and verification fails, the server will error out and refuse to
start. If set to false, a warning will be printed for any missing environment variables.

`--requireAllVars=true`

**Default:** `false`
