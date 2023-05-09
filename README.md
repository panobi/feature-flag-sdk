# Panobi Feature Flag SDK

## Overview

This SDK lets you push feature flag changes to your Panobi workspace, so that release statuses will update automatically and make it easier to interpret changes to your top-line metrics.

## Who is it for?

If you use an in-house feature flag system, you’re in the right place. This SDK doesn’t enable integration with third-party services like LaunchDarkly, Split, and Statsig (see Panobi’s custom integrations for those).

## How does it work?

The SDK is based on events. Each event represents changes to a flag, like the flag state — enabled or disabled — and the name of the flag. Events can be sent one at a time or in batches of up to 64.

## Compatibility

The [API specification](openapi.yaml) was generated for [OpenAPI 3.0.0](https://spec.openapis.org/oas/v3.0.0).

The source files were written against [Go 1.18](https://go.dev/doc/go1.18). They are known to work with [Go 1.20](https://go.dev/doc/go1.20), and may work with older versions.

## Getting started

The quickest way to get up and running is to run the provided [example programs](#running-the-example-programs), which demonstrate one method for constructing events and sending them to Panobi.

TODO: Integrate your in-house applications.

TODO: Use the provided OpenAPI specification. You can send events to Panobi over HTTP from the language or tool of your choice.

## Running the example programs

You will need your signing key, which you can copy from the integration settings in your Panobi workspace. The example programs expect the signing key in the form of an environment variable.

```console
export FEATURE_FLAG_SDK_SIGNING_KEY=<your signing key>
```

The simple example is a good place to start.

```console
cd examples/simple
go run main.go
```

Roughly, it works as follows.

1. Reads the enviroment variable and parses your key.
2. Creates a client with the parsed key.
3. Constructs an event representing an enabled feature flag.
4. Sends the event to Panobi.

Once the event has been successfully sent, it is available for use in your Panobi workspace. You should be able to select it from a drop-down menu in the editor panel for any Release. One selected, the state of the flag will be reflected in the Release. For example, if the flag is enabled, then the Release will be marked as Live, and moved into the appropriate column inside Panobi. If the flag is then  disabled via a subsequent event, the Release will be marked as Complete.

## License

## About Panobi

The platform designed for growth teams.

Panobi helps growth teams increase their velocity, deliver results, and amplify customer insights across the company.
