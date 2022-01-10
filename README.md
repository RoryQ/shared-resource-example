# Shared Resource Example

This example demonstrates sharing a single resource (e.g. a subscription) among multiple
kubernetes pods by using a distributed lock to prevent concurrent access. 

## Overview

This example uses a [distributed lock library](https://github.com/flowerinthenight/spindle)
built on Google Cloud Spanner as the locking mechanism. 
- `example/spanner` directory manages a spanner emulator pod and runs the migrations.
- `example/protected` directory contains a dummy service which we want to limit access
to a single process. It has endpoints to increment / decrement a counter
- `example/app` is the app that utilises the distributed lock to access the shared resource.

The app attempts to access the shared resource with the following method
1. Wait for a specific label/value to be set on the pod. i.e. to indicate when a blue/green deployment
has completed.
2. Attempt to take the distributed lock. When taken the library will maintain the lock using a background process
3. Use the shared resource by calling `/connect`, when done (or the service is being shut down) it will call `/disconnect`



## How to run

Requires
- skaffold
- minikube

Once installed and minikube is running then run skaffold to deploy:

```bash
skaffold run --tail --port-forward
```

## Links

