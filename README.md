# LaunchDarkly Go SDK with Fastly Compute

> ⚠️ **Warning**: This is an example implementation meant to serve as a starting point for integrating LaunchDarkly with Fastly Compute. The code is incomplete and should be used at your own risk. Make sure to:
>
> - Add proper error handling
> - Implement security best practices
> - Add appropriate logging
> - Test thoroughly in your environment
> - Follow LaunchDarkly and Fastly's latest best practices

This project demonstrates how to integrate the LaunchDarkly Go SDK with Fastly's Compute Runtime. It showcases a practical implementation of feature flagging in a serverless environment, specifically using Fastly's KV Store for persistent data storage.

## Overview

This demo application shows how to:

- Initialize the LaunchDarkly Go SDK in Fastly's Compute Runtime
- Use Fastly's KV Store as a persistent data store for LaunchDarkly
- Handle feature flag evaluations in a serverless context
- Maintain feature flag state across requests

## Key Components

### KV Store Implementation

The `kvdatastore` package provides a custom implementation of LaunchDarkly's data store interface using Fastly's KV Store. This allows feature flag data to persist between requests, improving performance and reducing API calls to LaunchDarkly's servers.

### Main Application

The main application (`main.go`) demonstrates:

- Initialization of the LaunchDarkly client with custom configuration
- Context creation with Fastly-specific attributes
- Feature flag evaluation
- JSON response formatting

## Prerequisites

- Go 1.21 or later
- Fastly CLI
- LaunchDarkly account and SDK key
- Fastly account with Compute@Edge enabled

## Configuration

1. Set up your LaunchDarkly SDK key in the environment or update the constant in `main.go`:

```go
const LD_SDK_KEY = "your-sdk-key"
```

2. Configure your Fastly service using the `fastly.toml` file.

## Local Development

1. Install dependencies:

```bash
go mod download
```

2. Run the application locally:

```bash
fastly compute serve
```

## Deployment

Deploy to Fastly using:

```bash
fastly compute publish
```

## How It Works

1. When a request comes in, the application initializes the LaunchDarkly client with custom configuration using Fastly's KV Store.
2. The client creates a context with Fastly-specific attributes (POP, region, service version, etc.).
3. Feature flags are evaluated using this context.
4. Results are returned as JSON responses.

## Response Format

The application returns JSON responses in the following format:

```json
{
  "animal": "feature-flag-value",
  "context": {
    "kind": "multi",
    "contexts": [
      {
        "kind": "user",
        "key": "user-123"
      },
      {
        "kind": "fastly-request",
        "key": "request-id",
        "fastly_service_version": "version",
        "fastly_pop": "pop",
        "fastly_region": "region",
        "fastly_service_id": "service-id"
      }
    ]
  },
  "reason": {
    "kind": "RULE_MATCH",
    "ruleIndex": 0,
    "ruleId": "rule-id"
  },
  "serviceVersion": "version"
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
