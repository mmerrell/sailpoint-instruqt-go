---
slug: local-activities
id: kby5u4sntxvp
type: challenge
title: 'Exercise 4: Local Activities'
teaser: Move fast in-process steps to local activities to reduce task queue round-trips.
notes:
- type: text
  contents: |-
    `ValidateOrder` and `FraudCheck` are pure in-process checks — no network calls,
    no external dependencies. Running them as regular activities wastes two task queue
    round-trips to the Temporal Server per workflow execution.

    At scale (20K+ orders per day), that is 40K unnecessary Actions.

    **Local activities eliminate the round-trip.** They run inside the worker process
    and produce a single `MarkerRecorded` event in the history instead of the
    schedule/start/complete triplet you get from regular activities.

    Hit **Start** when you're ready.
tabs:
- id: tckpi2a8mwtd
  title: Code
  type: code
  hostname: workshop
  path: /workspace/exercise
- id: imtw6lyxjpma
  title: Terminal 1 - Worker
  type: terminal
  hostname: workshop
  workdir: /workspace/exercise
- id: t11jlutq46l2
  title: Terminal 2 - Starter
  type: terminal
  hostname: workshop
  workdir: /workspace/exercise
- id: oocix3bslhma
  title: Temporal Web UI
  type: service
  hostname: workshop
  path: /
  port: 8080
difficulty: basic
timelimit: 3600
enhanced_loading: null
---

## Exercise 4: Local Activities

Open **`fulfillment_workflow.go`** and **`local_activities.go`** in the [button label="Code" background="#444CE7"](tab-0) tab.
The local activity implementations are already in `local_activities.go`.
Your job is to call them from the workflow using the local activity API.

***

### Part A – Set up local activity options

In `fulfillment_workflow.go`, add local activity options **before** the child workflow call:

```go
lao := workflow.LocalActivityOptions{
    StartToCloseTimeout: 5 * time.Second,
}
localCtx := workflow.WithLocalActivityOptions(ctx, lao)
lfa := &LocalFulfillmentActivities{}
```

> Use `workflow.LocalActivityOptions` — not `workflow.ActivityOptions`.
> Use `workflow.WithLocalActivityOptions` — not `workflow.WithActivityOptions`.

***

### Part B – Call `ValidateOrder`

`ValidateOrder` returns only an error — no result to capture:

```go
if err := workflow.ExecuteLocalActivity(localCtx, lfa.ValidateOrder, order).Get(localCtx, nil); err != nil {
    return OrderResult{}, err
}
```

***

### Part C – Call `FraudCheck`

`FraudCheck` returns a risk score string:

```go
var riskScore string
if err := workflow.ExecuteLocalActivity(localCtx, lfa.FraudCheck, order).Get(localCtx, &riskScore); err != nil {
    return OrderResult{}, err
}
log.Info("Fraud check passed", "riskScore", riskScore)
```

***

### Part D – Run it

1. Click the [button label="Terminal 1 - Worker" background="#444CE7"](tab-1) tab and start the Worker:

   ```bash,run
   go run ./cmd/worker/
   ```

2. Click the [button label="Terminal 2 - Starter" background="#444CE7"](tab-2) tab and run the Starter:

   ```bash,run
   go run ./cmd/starter/
   ```

In the [button label="Temporal Web UI" background="#444CE7"](tab-3) tab, open `fulfillment-ORD-4001`. Look for `MarkerRecorded` events
for the local activities — compare with the `ActivityTaskScheduled` /
`ActivityTaskStarted` / `ActivityTaskCompleted` triplet for the remote activities.

The two local activity calls produced **2 events** total. As regular activities
they would have produced **6 events**. At 20K workflows/day that is 80K fewer Actions.

***

Click **Check** when done, or **Solve** to see the reference solution.
