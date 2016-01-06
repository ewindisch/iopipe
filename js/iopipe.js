/*
 Usage:
 var iopipe = require('iopipe')
 iopipe.exec("http://api.twitter.com/blah/blah", function() {}, "sha256:DEADBEEF", "user/pipeline", "http://somedestination/")
*/
var url = require('url')
var request = require("request")
var util = require('util')
var vm = require('vm')
var fs = require('fs')

function funcCallback(call, done) {
  return function() {
    done(call.apply(this, [].slice.call(arguments)))
  }
}

function httpCallback(u, done) {
  return function() {
    if (arguments.length === 0) {
      request.get({url: url.format(u), strictSSL: true }, function(error, response, body) {
        done(body)
      })
    } else {
      prevResult = arguments[0]
      request.post({url: url.format(u), body: prevResult, strictSSL: true },
                    function(error, response, body) {
                      done(body)
                    })
    }
  }
}

function pipescriptCallback(id, done) {
  // Pull from index (or use cached pipescripts)
  /* download script */
  var script = fs.readFileSync(".iopipe/filter_cache/" + id)
  var input = ""
  var svm = vm.Script(script)

  return function(prevResult) {
    console.log("Running script")
    var prevResult = ""
    if (arguments.length > 0) {
      prevResult = arguments[0]
    }
    var result = svm.runInNewContext({ input: prevResult })
    console.log("Done: " + result)
    return done(result)
  }
}

exports.define = function() {
  var callbackList = []
  var nextCallback;
  var done = function(result) { console.log(result) };

  for (var i = arguments.length - 1; i > -1; i--) {
    var arg = arguments[i];

    if (typeof arg === "function") {
      nextCallback = funcCallback(arg, done)
    } else if (typeof(arg) === "string") {
      var u = url.parse(arg);

      if (u.protocol === 'http:' || u.protocol === 'https:') {
        var server = u.hostname
        nextCallback = httpCallback(u, done)
      } else {
        nextCallback = pipescriptCallback(arg, done)
      }
    } else {
      throw new Error("ERROR: unknown argument: " + arg)
    }
    done = nextCallback
  }
  return nextCallback
}

exports.exec = function() {
  var l = [].slice.call(arguments)
  exports.define.apply(this, l)()
}
