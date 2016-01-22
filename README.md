IOpipe
---------------------------------------
Apache 2.0 licensed.

IOpipe simplifies the consumption and integration of web services through
the chaining of pipescripts, javascript-based microservices.

Pipescripts take and transform input, providing straight-forward output
in a fashion to Unix pipes. Scripts may receive input or send output to/from
web service requests, functions, or local applications.

Scripts may be embedded in applications, used from shell scripts, or run manually
via a CLI. Because pipescript is Javascript, embedding is possible with most
languages and made safe through sandboxing.

![Build Status](https://circleci.com/gh/iopipe/iopipe.png?circle-token=eae431abda6b19dbfca597af818bb01092211272)

---------------------------------------
Usage
---------------------------------------

### Installation:

Download the [latest binary release](https://github.com/iopipe/iopipe/releases) and chmod 755 the file.

Building from source? See [Build & Install from source](#build--install-from-source).

Alternatively, download & alias our Docker image:

```bash
$ docker pull iopipe/iopipe:trunk
$ docker run --name iopipe-data iopipe/iopipe:trunk
$ eval $(echo "alias iopipe='docker run --rm --volumes-from iopipe-data iopipe/iopipe:trunk'" | tee -a ~/.bashrc)
$ iopipe --help
```

OS-specific packages are forthcoming.

### Command-line

```sh
# Import a pipescript and name it com.example.SomeScript
$ iopipe import --name com.example.SomeScript - <<<'input'

# List pipescripts
$ iopipe list

# Fetch response and process it with com.example.SomeScript
$ iopipe --debug exec http://localhost/some-request com.example.SomeScript

# Fetch response and convert it with SomeScript, sending the result to otherhost
$ iopipe --debug exec http://localhost/some-request com.example.SomeScript \
                      http://otherhost/request

# Fetch response and convert it with SomeScript, send that result to otherhost,
# & converting the response with the script ResponseScript
$ iopipe --debug exec http://localhost/some-request com.example.SomeScript \
                      http://otherhost/request some.example.ResponseScript
```

### NodeJS SDK:

The NodeJS SDK provides a generic callback chaining mechanism which allows
mixing HTTP(S) requests/POSTs, function calls, and pipescripts. Callbacks
receive the return of the previous function call or HTTP body.

```javascript
var iopipe = require("iopipe")

// Where com.example.SomeScript is present in .iopipe/filter_cache/
iopipe.exec("http://localhost/get-request",
            "com.example.SomeScript",
            "http://otherhost.post")

// Users may chain functions and HTTP requests.
iopipe.exec(function() { return "something" },
            function(arg) { return arg },
            "http://otherhost.post",
            your_callback)

// A function may also be returned then executed later.
var f = iopipe.define("http://fetch", "https://post")
f()

// A defined function also accepts parameters
var echo = require("iopipe-echo")
var f = iopipe.define(echo, console.log)
f("hello world")
```

For more information on using the NodeJS SDK, please refer to its documentation:
***https://github.com/iopipe/iopipe/blob/master/docs/nodejs.md***

---------------------------------------
Filters & Pipescript
---------------------------------------

Requests and responses and translated using filters written in
Pipescript (i.e. Javascript) or offered as web services.

All filters simply receive request or response data and output
translated request or response data. Pipescript is typically operated
upon locally in the client, whereas web-service based filters operate
server-side. Pipescript may also be used to build serverside filters
and applications.

Example:

```javascript
module.exports = function(input) {
  return "I'm doing something with input: {0}".format(input)
}
```

For more on writing filters see:
***https://github.com/iopipe/iopipe/blob/master/docs/pipescript.md***

---------------------------------------
Build & Install from source
---------------------------------------

With a functioning golang 1.5 development environment:

```bash
$ go build
$ ./iopipe --help
```

Alternatively use Docker to build & deploy:

```bash
$ docker build -t iopipe-dev .
$ docker run --name iopipe-data iopipe-dev
$ eval $(echo "alias iopipe='docker run --rm --volumes-from iopipe-data iopipe-dev'" | tee -a ~/.bashrc)
$ iopipe --help
```

---------------------------------------
Project goals
---------------------------------------

The principal goal of our project is to improve
human to machine and machine to machine communication.
We believe this can be achieved without the creation
or use of new protocols, but through the use of
a flow-based programming model.

Furthermore:

1. Simplify the use and integration of existing APIs into
   user applications.
2. Support use by both existing and new applications.
3. Design for an open and distributed web.
4. Permissive open source licensing.
5. Secure sharing, execution, & communications.

---------------------------------------
Scaling
---------------------------------------

Filters are pulled from (de)centralized repositories
and scale should be considered when building and
deploying private filter repositories.

---------------------------------------
Security
---------------------------------------

Note that this tool communicates to 3rd-party
web services. Caution is advised when trusting
these services, as is standard practice with
web and cloud services.

Pipescripts are executed in individual
javascript VMs whenever allowed by the executing
environment. Users should exercise caution
when running community created pipescripts.

It is a project priority to make fetching, publishing,
and execution of pipescripts secure for a
production-ready 1.0.0 release.

Contact security@iopipe.com for questions.

---------------------------------------
LICENSE
---------------------------------------

Apache 2.0
