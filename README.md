# Panobi Feature Flag SDK

## Overview

This SDK lets you push feature flag changes from your custom project management system to your Panobi workspace, so that Panobi project statuses will update automatically and make it easier to interpret changes to your top-line metrics.

## Who is it for?

If you use an in-house feature flag system, you’re in the right place. This SDK doesn’t enable integration with third-party services like LaunchDarkly, Split, and Statsig (see Panobi’s custom integrations for those).

## How does it work?

The SDK is based on events. Each event represents changes to a flag, like the flag state — enabled or disabled — and the name of the flag. Events can be sent one at a time or in batches of up to 64.

## Compatibility

The [API specification](openapi.yaml) was generated for [OpenAPI 3.0.0](https://spec.openapis.org/oas/v3.0.0).

## Getting started

You will need your signing key, which you can copy from the integration settings in your Panobi workspace. Your key has three parts, separated by dashes, the last of which is _secret_. You should never share this secret publicly. Nor should you pass it directly to Panobi when making requests. Rather, it should be used to calculate a signature for your request, the process for which is illustrated in a short script further down in this document.

To send an event to your Panobi workspace, make an HTTP POST request with a JSON body to the following endpoint:

https://app.panobi.com/integrations/flags-sdk/events

The body of the request should be JSON encoded. The complete format can be found in the provided [OpenAPI specification](openapi.yaml). For example, the following payload would mark a single flag named `Beta Feature XYZ` as live:

```json
{
    "events": [
        {
            "project": "growth-team",
            "key": "beta-feature-xyz",
            "name": "Beta Feature XYZ",
            "dateModified": "2023-03-10T17:27:55+00:00",
            "isEnabled": true
        }
    ]
}
```

### Request signing

Once you've built your request, you need to sign it so that Panobi knows it's from you. This is where the secret part of your signing key comes into play.

The little shell script below demonstrates the complete process, including how to split your signing key, how to generate a timestamp (to prevent replay attacks) and a signature, and which headers to send us.

```shell
#!/usr/bin/env bash

# We read the request body from a file, so that we can generate a signature for
# it. The first argument to this script is the name of that file.
input=$(<"$1")

# Check that we have a signing key
if [[ -z "$FEATURE_FLAG_SDK_SIGNING_KEY" ]]; then
	echo "Did you forget to set \$FEATURE_FLAG_SDK_SIGNING_KEY?"
	exit 1
fi

# Split the signing key into its component parts.
arr=(${FEATURE_FLAG_SDK_SIGNING_KEY//-/ })
wid=${arr[0]}    # Workspace ID
eid=${arr[1]}    # External ID
secret=${arr[2]} # Secret

# Get the milliseconds since Unix epoch. We'll use this as a timestamp for the
# signature to prevent replay attacks.
ts=$(date +%s)000

# Hash the timestamp and the request body using the secret part of the signing
# key.
msg="v0:${ts}:${input}"
sig=$(echo -n "${msg}" | openssl dgst -r -sha256 -hmac "${secret}" | awk '{print $1}')

# Post the headers and the request to Panobi using good ol' curl.
curl -v \
    -X POST \
    -H "X-Panobi-Signature: v0=""${sig}" \
    -H "X-Panobi-Request-Timestamp: ""${ts}" \
    -H "Content-Type: application/json" \
    -d "${input}" \
    https://app.panobi.com/integrations/flags-sdk/events/"${wid}"/"${eid}"
```

Once the event has been successfully sent, the flag is available for use in your Panobi workspace. You should be able to select it from a drop-down menu in the editor panel for any Panobi Project. Once selected, the state of the flag will be reflected in the Panobi Project. For example, if the flag is enabled, then the Panobi Project will be marked as Live, and moved into the appropriate column inside Panobi. If the flag is then disabled via a subsequent event, the Panobi Project will be marked as Complete.

## License

This SDK is provided under the terms of the [Apache License 2.0](LICENSE).

## About Panobi

Panobi is the platform for growth observability: helping companies see, understand, and drive their growth.
