# Effe, a simple Open Source building block to emulate AWS Lambda

Effe, is an etremely simple building block to build a "server-less" architecture.

This is a building block, it manage a very single lambda function.

## Terminology

*Lambda*: a single, simple function that speak HTTP, you can run multiple instance of a single lambda at once, an example of lambda is a funcion of AWS Lambda.

*effe*: An `effe` is a function 

## Life cicle

To use `effe` you can provide four different functions,`logic.Init()`, `logic.Start() Context`, `logic.Run(Context, http.ResponseWritter, *http.Request) error` and `logic.Stop(Context)` and a custom type: Context.

#### Init

The function `logic.Init()` it is run once and only once when an effe is started, you can use it to set the seed of the random number generator, communicate with a discovery service, or whatever.

#### Start

The function `logic.Start() Context` takes no argument and return a `Context`.

This function it can be run multiple times if the effe is under heavy load, or just once, if the effe is not under load.

You can initialize connections to databases, communicate with a discovery services, or pretty much whatever.

You can pack whatever information inside the Context that it will be passed to `logic.Run`

#### Run

This function is called every time an effe is invoke, it is pretty similar to a simple `http handler` the only difference is that is take also a `Context` as first argument.

Yourself have compose the `Context` in the `Start` function so it is very powerfull, however you definitely shouldn't save any form of state inside the `Context`.

This function can read the request and provide a responses.

It can also return an `error`, the error will not be handle but it will be logged to syslog, if everything went good and nothing need to be logged you should return `nil`.

#### Stop

This function takes the `Context` and gracefully destroy it. 
If you have open a database connection you may want to close it.

#### Context




## How to use

My idea is to run one or more docker containers for every lambda.

However our effes doesn't know what resource/URL it should respond, and we like this, but still we need a way to route the traffic.

To route the traffic we can write a proxy server that simply forward the calls which is a simple to solution but may have scaling problem and it is somehow slow. 
There is the need for the proxy to accept the external call, decide what effe is necessary to call, make another HTTP request, wait on the result, and finaly send the result back.
In this case however you only need one single IP address exposed, the effes can run inside not exposed container.

Another way to route the traffic is to use HAproxy which doesn't accept the HTTP request but simply re-route it to a particular server. In this case there is only one HTTP request instead of two but every single effe need to be exposed to the web.

Either way you will need to manage a lot of containers, it is a good idea to use an automatic tool to create the docker instances, I would suggest to look into [Kubernetes](kubernetes) which already provide a way to route the traffic via [Ingress](ingress).

## Docker 

Effe can be compiled down to sigle binary, it means that you can have a very light docker images with inside only the necessary logic to run effe.

However since we want effe to use http*s* is necessary to provide the certificates, it means that we cannot use the `SCRATCH` images from docker but we need to include some certificates. A docker container with nothing but the certificates is [centurylink/ca-certs](ca-certs).

Of course you can feel free to build a similar images with only the certificates you trust and maybe you could also share such image.

Surely you can use your own container with whatever images, but it is probably going to be bigger than this strategy.


[kubernetes]: http://kubernetes.io/
[ingress]: http://kubernetes.io/v1.1/docs/user-guide/ingress.html
[ca-certs]: https://hub.docker.com/r/centurylink/ca-certs/
