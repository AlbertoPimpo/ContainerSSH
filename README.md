[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">An SSH Server that Launches Containers in Kubernetes and Docker</h1>

[![Documentation: available](https://img.shields.io/badge/documentation-available-green)](https://containerssh.io/)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/ContainerSSH/ContainerSSH/main.yml?&branch=main)](https://github.com/containerssh/containerssh/actions)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/containerssh/containerssh?sort=semver)](https://github.com/containerssh/containerssh/releases)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/containerssh/containerssh)](http://hub.docker.com/r/containerssh/containerssh)
[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/containerssh)](https://goreportcard.com/report/github.com/containerssh/containerssh)
[![License: Apache 2.0](https://img.shields.io/github/license/ContainerSSH/ContainerSSH)](LICENSE.md)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FContainerSSH%2FContainerSSH.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2FContainerSSH%2FContainerSSH?ref=badge_shield&issueType=license)

## ContainerSSH in One Minute

In a hurry? This one-minute video explains everything you need to know about ContainerSSH.

[![An image with a YouTube play button on it.](https://containerssh.io/images/containerssh-intro-preview.png)](https://youtu.be/Cs9OrnPi2IM)

## Need help?

[Join the #containerssh Slack channel on the CNCF Slack »](https://communityinviter.com/apps/cloud-native/cncf)

## Use cases

### Build a lab

Building a lab environment can be time-consuming. ContainerSSH solves this by providing dynamic SSH access with APIs, automatic cleanup on logout using ephemeral containers, and persistent volumes for storing data. **Perfect for vendor and student labs.**

[Read more »](https://containerssh.io/usecases/lab/)

### Debug a production system

Provide **production access to your developers**, give them their usual tools while logging all changes. Authorize their access and create short-lived credentials for the database using simple webhooks. Clean up the environment on disconnect.

[Read more »](https://containerssh.io/usecases/debugging/)

### Run a honeypot

Study SSH attack patterns up close. Drop attackers safely into network-isolated containers or even virtual machines, and **capture their every move** using the audit logging ContainerSSH provides. The built-in S3 upload ensures you don't lose your data.

[Read more »](https://containerssh.io/usecases/honeypots/)

## How does it work?

![](https://containerssh.io/images/architecture.svg)

1. The user opens an SSH connection to ContainerSSH.
2. ContainerSSH calls the authentication server with the users username and password/pubkey to check if its valid.
3. ContainerSSH calls the config server to obtain backend location and configuration (if configured)
4. ContainerSSH calls the container backend to launch the container with the
   specified configuration. All input from the user is sent directly to the backend, output from the container is sent
   to the user.

[▶️ Watch as video »](https://youtu.be/Cs9OrnPi2IM) | [🚀 Get started »](https://containerssh.io/quickstart/)

## Demo

![](https://containerssh.io/images/ssh-in-action.gif)

[🚀 Get started »](https://containerssh.io/quickstart/)

## Verify provenance

Each of the releases come with a SLSA provenance data file `multiple.intoto.jsonl`. This file can be used to verify the source and provenance of the produced artifacts with [`slsa-verifier`](https://github.com/slsa-framework/slsa-verifier).


This aims to ensure the users that the artifacts are coming from containerssh.

An example of verification :
```sh
slsa-verifier verify-artifact <artifact-to-verify> \
--provenance-path <path-to-your-provenance> \
--source-uri github.com/containerssh/containerssh
```

If the verification is successful, the process should produce the following output :
```
Verifying artifact <artifact-to-verify>: PASSED
PASSED: Verified SLSA provenance
```


## Contributing

If you would like to contribute, please check out our [Code of Conduct](https://github.com/ContainerSSH/community/blob/main/CODE_OF_CONDUCT.md) as well as our [contribution documentation](https://containerssh.io/development/).

## Embedding ContainerSSH

You can fully embed ContainerSSH into your own application. First, you will need to create the configuration structure:

```go
cfg := config.AppConfig{}
// Set the default configuration:
cfg.Default()
```

You can then populate this config with your options and create a ContainerSSH instance like this:

```go
pool, lifecycle, err := containerssh.New(cfg, loggerFactory)
if err != nil {
    return err
}
```

You will receive a service pool and a lifecycle as a response. You can use these to start the service pool of ContainerSSH. This will block execution until ContainerSSH stops.

```go
err := lifecycle.Run()
```

This will run ContainerSSH in the current Goroutine. You can also use the lifecycle to add hooks to lifecycle states of ContainerSSH. You must do this *before* you call `Run()`. For example:

```go
lifecycle.OnStarting(
    func(s service.Service, l service.Lifecycle) {
        print("ContainerSSH is starting...")
    },
)
```

You can also have ContainerSSH stop gracefully by using the `Stop()` function on the lifecycle. This takes a context as an argument, which is taken as a timeout for the graceful shutdown.

Finally, you can use the returned `pool` variable to rotate the logs. This will trigger all ContainerSSH services to close and reopen their log files.

```
pool.RotateLogs()
```

## Building an authentication webhook server

## Building a configuration webhook server

The configuration webhook lets you dynamically configure ContainerSSH. This library contains the tools to create a tiny webserver to serve these webhook requests.

First, you need to fetch this library as a dependency using [go modules](https://blog.golang.org/using-go-modules):

```bash
go get go.containerssh.io/containerssh
```

Next, you will have to write an implementation for the following interface:

```go
package main

import (
    "go.containerssh.io/containerssh/config"
)

type ConfigRequestHandler interface {
	OnConfig(request config.Request) (config.AppConfig, error)
}
```

The best way to do this is creating a struct and adding a method with a receiver:

```go
type myConfigReqHandler struct {
}

func (m *myConfigReqHandler) OnConfig(
    request configuration.ConfigRequest,
) (config configuration.AppConfig, err error) {
    // We recommend using an IDE to discover the possible options here.
    if request.Username == "foo" {
        config.Docker.Config.ContainerConfig.Image = "yourcompany/yourimage"
    }
    return config, err
}
```

**Warning!** Your `OnConfig` method should *only* return an error if it can genuinely not serve the request. This should not be used as a means to reject users. This should be done using the authentication server. If you return an error ContainerSSH will retry the request several times in an attempt to work around network failures.

Once you have your handler implemented you must decide which method you want to use for integration.

### The full server method

This method is useful if you don't want to run anything else on the webserver, only the config endpoint. You can create a new server like this:

```go
package main

import (
	"signal"
	
	"go.containerssh.io/containerssh/config"
	"go.containerssh.io/containerssh/config/webhook"
	"go.containerssh.io/containerssh/log"
	"go.containerssh.io/containerssh/service"
)

func main() {
	logger := log.NewLogger(&config.LogConfig{
		// Add logging configuration here
    })
	// Create the webserver service
    srv, err := webhook.NewServer(
        config.HTTPServerConfiguration{
            Listen: "0.0.0.0:8080",
        },
        &myConfigReqHandler{},
        logger,
    )
	if err != nil {
		panic(err)
    }

	// Set up the lifecycle handler
	lifecycle := service.NewLifecycle(srv)
	
	// Launch the webserver in the background
	go func() {
		//Ignore error, handled later.
		_ = lifecycle.Run()
	}()

    // Handle signals and terminate webserver gracefully when needed.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if _, ok := <-signals; ok {
			// ok means the channel wasn't closed, let's trigger a shutdown.
			// The context given is the timeout for the shutdown.
			lifecycle.Stop(
				context.WithTimeout(
					context.Background(),
					20 * time.Second,
				),
			)
		}
	}()
	// Wait for the service to terminate.
	lastError := lifecycle.Wait()
	// We are already shutting down, ignore further signals
	signal.Ignore(syscall.SIGINT, syscall.SIGTERM)
	// close signals channel so the signal handler gets terminated
	close(signals)

	if lastError != nil {
		// Exit with a non-zero signal
		fmt.Fprintf(
			os.Stderr,
			"an error happened while running the server (%v)",
			lastError,
		)
		os.Exit(1)
	}
	os.Exit(0)
}
```

**Note:** We recommend securing client-server communication with certificates. Please see the [Securing webhooks section below](#securing-webhooks)/

### Integrating with an existing HTTP server

Use this method if you want to integrate your handler with an existing Go HTTP server. This is rather simple:

```go
handler, err := configuration.NewHandler(&myConfigReqHandler{}, logger)
```

You can now use the `handler` variable as a handler for the [`http` package](https://golang.org/pkg/net/http/) or a MUX like [gorilla/mux](https://github.com/gorilla/mux).

## Using the config client

This library also contains the components to call the configuration server in a simplified fashion. To create a client simply call the following method:

```go
client, err := configuration.NewClient(
	configuration.ClientConfig{
        http.ClientConfiguration{
            URL: "http://your-server/config-endpoint/"
        }
    },
	logger,
    metricsCollector,
)
```

The `logger` is a logger from the [log library](https://github.com/ContainerSSH/ContainerSSH/tree/main/log), the `metricsCollector` is supplied by the [metrics library](https://github.com/containerssh/tree/main/metrics).

You can now use the `client` variable to fetch the configuration specific to a connecting client:

```go
connectionID := "0123456789ABCDEF"
appConfig, err := client.Get(
    ctx,
    "my-name-is-trinity",
    net.TCPAddr{
        IP: net.ParseIP("127.0.0.1"),
        Port: 2222,
    },
    connectionID,
) (AppConfig, error)
```

Now you have the client-specific configuration in `appConfig`.

**Note:** We recommend securing client-server communication with certificates. The details about securing your HTTP requests are documented in the [HTTP library](https://github.com/containerssh/containerssh/tree/main/http).

## Loading the configuration from a file

This library also provides simplified methods for reading the configuration from an `io.Reader` and writing it to an `io.Writer`.

```go
file, err := os.Open("file.yaml")
// ...
loader, err := configuration.NewReaderLoader(
	file,
    logger,
    configuration.FormatYAML,
)
// Read global config
appConfig := &configuration.AppConfig{}
err := loader.Load(ctx, appConfig)
// Read connection-specific config:
err := loader.LoadConnection(
    ctx,
    "my-name-is-trinity",
    net.TCPAddr{
        IP: net.ParseIP("127.0.0.1"),
        Port: 2222,
    },
    connectionID,
    appConfig,
)
```

As you can see these loaders are designed to be chained together. For example, you could add an HTTP loader after the file loader:

```go
httpLoader, err := configuration.NewHTTPLoader(clientConfig, logger)
```

This HTTP loader calls the HTTP client described above.

Conversely, you can write the configuration to a YAML format:

```go
saver, err := configuration.NewWriterSaver(
    os.Stdout,
    logger,
    configuration.FormatYAML,
)
err := saver.Save(appConfig)
```


## Building a combined configuration-authentication webhook server

## Securing webhooks

## Reading audit logs
