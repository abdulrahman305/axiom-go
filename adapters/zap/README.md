# Axiom Go Adapter for uber-go/zap

Adapter to ship logs generated by [uber-go/zap](https://github.com/uber-go/zap)
to Axiom.

## Quickstart

Follow the [Axiom Go Quickstart](https://github.com/axiomhq/axiom-go#quickstart)
to install the Axiom Go package and configure your environment.

Import the package:

```go
// Imported as "adapter" to not conflict with the "uber-go/zap" package.
import adapter "github.com/axiomhq/axiom-go/adapters/zap"
```

You can also configure the adapter using [options](https://pkg.go.dev/github.com/axiomhq/axiom-go/adapters/zap#Option)
passed to the [New](https://pkg.go.dev/github.com/axiomhq/axiom-go/adapters/zap#New)
function:

```go
core, err := adapter.New(
    SetDataset("AXIOM_DATASET"),
)
```

To configure the underlying client manually either pass in a client that was
created according to the [Axiom Go Quickstart](https://github.com/axiomhq/axiom-go#quickstart)
using [SetClient](https://pkg.go.dev/github.com/axiomhq/axiom-go/adapters/zap#SetClient)
or pass [client options](https://pkg.go.dev/github.com/axiomhq/axiom-go/axiom#Option)
to the adapter using [SetClientOptions](https://pkg.go.dev/github.com/axiomhq/axiom-go/adapters/zap#SetClientOptions).

```go
import adapter "github.com/axiomhq/axiom-go/axiom"

// ...

core, err := adapter.New(
    SetClientOptions(
        axiom.SetPersonalTokenConfig("AXIOM_TOKEN", "AXIOM_ORG_ID"),
    ),
)
```

### ❗ Important ❗

The adapter uses a buffer to batch events before sending them to Axiom. This
buffer must be flushed explicitly by calling [Sync](https://pkg.go.dev/github.com/axiomhq/axiom-go/adapters/zap#WriteSyncer.Sync). Refer to the
[zap documentation](https://pkg.go.dev/go.uber.org/zap/zapcore#WriteSyncer)
for details and checkout out the [example](../../examples/zap/main.go).