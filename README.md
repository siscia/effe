[![Build Status](https://travis-ci.org/siscia/effe.svg?branch=master)](https://travis-ci.org/siscia/effe)

# Effe, a simple Open Source building block to emulate AWS Lambda

Effe is an extremely simple building block with which to build a "server-less" architecture.

This is a building block, operates on the level of a single lambda function.


## Terminology

**Lambda**: a single, simple function that speaks HTTP.  You can run multiple instances of a single lambda at once. An example of a lambda is what is called a "function" in AWS Lambda.

**effe**: An `effe` is one of these little programs, compiled and running.

## Life cycle

To use `effe`, you can provide a custom type: `Context` and four different functions, `logic.Init()`, `logic.Start() Context`, `logic.Run(Context, http.ResponseWritter, *http.Request) error`, `logic.Stop(Context)`, and string variable `Info`.

A very basic example is available in `logic/logic.com`, which is also the fastest way to get started.

#### Init

The function `logic.Init()` it is run once and only once when an effe is started, you can use it to set the seed of the random number generator, communicate with a discovery service, or whatever.

#### Start

The function `logic.Start() Context` takes no argument and return a `Context` and an error.

This function can be run multiple times if the effe is under heavy load, or just once, if the effe is not under load.

You can initialize connections to databases, communicate with a discovery services, or pretty much whatever.

You can pack whatever information inside the Context and it will be passed to `logic.Run`

If everything goes well, the error should be `nil`, if some problems arise you can provide any error you want.

Return a not-nil error means that the context won't be put in cache, however, the computation that required the context will actually use it.

#### Run

This function is called every time an effe is invoked.  It is pretty similar to a simple `http handler`, the only difference is that it also takes a `Context` as the first argument and an `error` as second argument.

It is possible that some `Run` will runs with an error so, if it is possible for `Start` to return an error, then `Run` should be able to handle it.

You must compose the `Context` in the `Start` function.  So, it is very powerful, however you definitely shouldn't save any form of state inside the `Context`.

This function can read the request and provide a responses.

It can also return an `error`, the error will not be handled, but it will be logged to syslog, if everything went well and nothing needed to be logged you should return `nil`.

#### Stop

This function takes the `Context` and gracefully destroys it.
If you have open a database connection you may want to close it.

#### Context

Context is the type that you need to define with everything you need to run your effe.

Contexts can be created -- and as well destroyed -- at any time.  You should not save any state in a context since there are not any guarantees.

#### Info

Info is a JSON string.  This string is printed when your effe is invoked with the `-info` parameter.

All the information saved here can be used in a later stage of the process (eg. During deployment).

In particular the information associated with the key `name` and the key `version` are used by `effe-tool` to name the executable in a coherent way.

## Dependencies

`effe` dependencies are managed in the same way as the dependencies of any golang project.

As long as the dependencies are on your gopath, effe will compile.

## Getting start

The easiest and fastest way to use `effes` is using [effe-tool][effe-tool].

Effe-tool provides an automatic way to generate a new empty effe that you can populate with your logic, and a simple way to compile it.

If you don't want to use `effe-tool`, you can also clone this repo, modify the `logic.go` file in such a way that is coherent with your goals and run `go compile effe`.

You can run the binary with `./effe` and it should be exposed on port 8080 of your localhost.

If you prefer to expose some other port it is enough to launch `./effe -port $PORT`.  As an example, to listen on port 7050 you will launch `./effe -port 7050`.

It is also possible to compile everything down to a single executable.  Just run `go compile -buildmode=exe effe`.  The executable will be bigger, but it can run in a very small docker image, see more on this below.

If you are having problems running `effe` with `effe-tool` I made a short video, please note that the video refer to an older version of effe which required the `Start` function to return only a `Context`, the newer version require the `Start` function to return a `Context` **and** an `error`.

[![Effe Short Intro](http://www.youtube.com/watch?v=dItO9E29WRQ/0.jpg)](http://www.youtube.com/watch?v=dItO9E29WRQ)

## How to use

The idea is to run one -- or multiple -- docker containers for every lambda.

However, our effes don't know what resource/URL they should respond to, and we like this, but still we need a way to route the traffic.

To route the traffic we can write a proxy server that simply forwards the calls to the appropriate server.  This is a simple to solution, but may have scaling problems and it is somehow slow.
The proxy will need to accept the external call, decide which effe is need be called, make another HTTP request, wait on the result, and finally send the result back.
In this case however you only need a single IP address exposed.  The effes can run inside, rather than having a publically exposed container.

Another way to route the traffic is to use [HAProxy][haproxy] which doesn't accept the HTTP request but simply re-routes it to a particular server. In this case there is only one HTTP request instead of two, but every single effe needs to be exposed to the web.

Either way, you will need to manage a lot of containers, so it is a good idea to use an automatic tool to create the docker instances.  I would suggest looking into [Kubernetes][kubernetes], which already provides a way to route the traffic via [Ingress][ingress] using HAProxy.

## Docker

Effe can be compiled down to single binary. This means that you can have very light docker images with containing only the necessary logic to run effe.

However, since we want effe to use http**s**, it is necessary to provide the certificates.  This means that we cannot use the `SCRATCH` images from docker, but we need to include some certificates. A docker container with nothing but the certificates is [centurylink/ca-certs][ca-certs].

Of course you can feel free to build similar images with only the certificates you trust (and maybe you could also share such images).

Surely, you can use your own container with whatever images, but it is probably going to be bigger than this strategy.



[effe-tool]: https://github.com/siscia/effe-tool
[kubernetes]: http://kubernetes.io/
[ingress]: http://kubernetes.io/v1.1/docs/user-guide/ingress.html
[ca-certs]: https://hub.docker.com/r/centurylink/ca-certs/
[haproxy]: http://www.haproxy.org/
