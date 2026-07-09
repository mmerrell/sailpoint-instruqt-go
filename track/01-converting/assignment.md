---
slug: converting
id: tnpcxrrbxa8k
type: challenge
title: 'Exercise 1: Converting a Workflow'
teaser: Replace a fragile retry loop with Temporal Activities and a durable Workflow.
notes:
- type: text
  contents: |-
    The starting point for this exercise is `pipeline.go` — a deliberately fragile
    implementation using manual retry loops and local variable state.

    If the process dies mid-execution, the order is lost. If a step fails after five
    attempts, the whole pipeline crashes. Retry behavior is hardcoded and invisible.

    **Temporal fixes all of this.** The activities are already implemented in
    `activities.go`. Your job is to wire them together in `workflow.go`.

    Hit **Start** when you're ready.
tabs:
- id: crr6c0zsxnef
  title: Code
  type: code
  hostname: workshop
  path: /workspace/exercise
- id: jku51xpzcbvb
  title: Terminal 1 - Worker
  type: terminal
  hostname: workshop
  workdir: /workspace/exercise
- id: xzxd14dyoqyy
  title: Terminal 2 - Starter
  type: terminal
  hostname: workshop
  workdir: /workspace/exercise
- id: qkcpgjlwvuaf
  title: Temporal Web UI
  type: service
  hostname: workshop
  path: /
  port: 8080
difficulty: basic
timelimit: 2400
enhanced_loading: null
---

## Exercise 1: Converting a Workflow

Open **`pipeline.go`** and **`workflow.go`** in the [button label="Code" background="#444CE7"](tab-0) tab. `pipeline.go` shows the
fragile "before Temporal" implementation — read through it to understand what
each step does. Then open `workflow.go` and implement the TODOs.

The activities are already implemented for you in `activities.go`.

***

### Part A – Call `ReserveInventory`

In `workflow.go`, replace the first TODO:

```go
var reservationID string
if err := workflow.ExecuteActivity(ctx, fa.ReserveInventory, order).Get(ctx, &reservationID); err != nil {
    return OrderResult{}, err
}
```

> **Note:** `.Get(ctx, &reservationID)` is Go's equivalent of Python's `await` —
> it blocks until the activity completes and writes the result into `reservationID`.

***

### Part B – Call `ProcessPayment`

```go
var paymentConfirmation string
if err := workflow.ExecuteActivity(ctx, fa.ProcessPayment, order).Get(ctx, &paymentConfirmation); err != nil {
    return OrderResult{}, err
}
```

***

### Part C – Call `DispatchToFulfillment` and return the result

```go
var trackingNumber string
if err := workflow.ExecuteActivity(ctx, fa.DispatchToFulfillment, order, reservationID).Get(ctx, &trackingNumber); err != nil {
    return OrderResult{}, err
}

return OrderResult{
    OrderID:             order.OrderID,
    Status:              "FULFILLED",
    ReservationID:       reservationID,
    PaymentConfirmation: paymentConfirmation,
    TrackingNumber:      trackingNumber,
}, nil
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

Open the [button label="Temporal Web UI" background="#444CE7"](tab-3) tab and find `fulfillment-ORD-1001`. Watch the
activity retries happen automatically — no retry loops in your code.

Try killing the worker mid-execution (Ctrl+C in Terminal 1) and restarting it.
What happens?

***

Click **Check** when done, or **Solve** to see the reference solution.
