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

The quickest way to get up and running is to run the provided [example programs](#running-the-example-programs), which demonstrate how to construct events and send them to Panobi.

If you’re running into an alert in Panobi that no flags can be found, [try running the CSV program](#csv); this will help you push your feature flags into Panobi before the regularly scheduled push associated with your SDK.

If you're using a language other than Golang, or you'd rather roll-your-own, then take a look at how to send flags to us via [OpenAPI](#openapi).

## Running the example programs

You will need your signing key, which you can copy from the integration settings in your Panobi workspace. The example programs expect the signing key in the form of an environment variable.

```console
export FEATURE_FLAG_SDK_SIGNING_KEY=<your signing key>
```

### Simple

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

### CSV

This example program demonstrates how to send more than one event at a time. It will read a file of comma-separated values, where each line represents one event. If your feature flag system offers an export to CSV, then this is a great way to populate those flags in your Panobi workspace.

```console
cd examples/csv
go run main.go
```

Each row is in the following format:

```
Project, Key, DateModified, IsEnabled, Name
```

Where the last two columns are optional and can be omitted.

The following are all examples of valid rows:

```
growth-team,beta-feature-xyz,2023-03-10T17:27:55+00:00,true,Beta Feature XYZ
growth-team,beta-feature-abc,2023-03-10T17:30:55+00:00
```

### JSON

This example program works like the CSV example, but reads events in JSON format.

```console
cd examples/json
go run main.go
```

The following is an example of a valid row:

```json
{"project": "growth-team", "key": "beta-feature-xyz", "name": "Beta Feature XYZ", "dateModified": "2023-03-10T17:27:55+00:00", "isEnabled": true}
```

## OpenAPI

In an effort to be language agnostic, we've provided an [OpenAPI specification](openapi.yaml) that you can use to push events directly to Panobi.

Once you've built a request according to the specification, you need to sign it so that Panobi knows it's from you. The following little shell script demonstrates how to do this.

```shell
#!/usr/bin/env bash

# Assume the request body is in a file, and the filename is the first argument to this script.
input=$(<"$1")

# Now split the signing key into its component parts.
arr=(${FEATURE_FLAG_SDK_SIGNING_KEY//-/ })
wid=${arr[0]} # Workspace ID
eid=${arr[1]} # External ID
secret=${arr[2]} # Secret

# Get the unix epoch in milliseconds. We'll use this as a timestamp for the signature.
ts=$(date +%s)000

# Hash the timestamp and the request using the secret part of your signing key.
msg="v0:${ts}:${input}"
sig=$(echo -n "${msg}" | openssl dgst -sha256 -hmac "${secret}")

# Post the headers and the request to Panobi using CURL.
curl -v \
    -X POST \
    -H "X-Panobi-Signature: v0=""${sig}" \
    -H "X-Panobi-Request-Timestamp: ""${ts}" \
    -H "Content-Type: application/json" \
    -d "${input}" \
    https://panobi.com/integrations/flags-sdk/events/"${wid}"/"${eid}"
```

## License

## About Panobi

The platform designed for growth teams.

Panobi helps growth teams increase their velocity, deliver results, and amplify customer insights across the company.
