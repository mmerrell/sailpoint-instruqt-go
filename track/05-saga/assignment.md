---
slug: saga
type: challenge
title: 'Exercise 5: The Saga Pattern'
teaser: Implement compensating transactions so multi-step workflows clean up after
  themselves on failure.
notes:
- type: text
  contents: |-
    Distributed systems can't use a single database transaction to span multiple
    services. When your workflow calls three external APIs in sequence — reserve
    inventory, charge payment, dispatch shipment — and the third one fails, the
    first two have already committed.

    The **Saga pattern** solves this with **compensating transactions**: each
    forward step has a corresponding undo step. If the pipeline breaks, the
    workflow runs the compensations in reverse order, restoring consistent state.

    In the Go SDK this is a `defer` block playing the role Java gives to
    `try/catch`. Declare your state variables (`reservationID`,
    `paymentConfirmation`) as named results or before the `defer`, so the
    deferred func can see them. It only takes action if the workflow is
    returning with a non-nil error, and it only compensates for the steps that
    actually completed — guard each call on an empty-string check.

    Because workflow execution is durable, compensations will run to completion
    even if the Worker restarts mid-compensation.

    Hit **Start** when you're ready.
tabs:
- title: Code
  type: code
  hostname: workshop
  path: /workspace/exercise
- title: Terminal 1 - Worker
  type: terminal
  hostname: workshop
  workdir: /workspace/exercise
- title: Terminal 2 - Starter
  type: terminal
  hostname: workshop
  workdir: /workspace/exercise
- title: Temporal Web UI
  type: service
  hostname: workshop
  path: /
  port: 8080
difficulty: intermediate
timelimit: 2700
enhanced_loading: null
---

## Exercise 5: The Saga Pattern

All your work is in **`fulfillment_workflow.go`** (open it in the [button label="Code" background="#444CE7"](tab-0) tab).

Files are in `/workspace/exercise/`. Look for the two `// TODO` blocks.

***

### Part A – Understand the domain

Open `fulfillment_activities.go` in the same tab. Notice it has five methods:

**Forward steps** (run in order):
1. `ReserveInventory` — holds stock for the order
2. `ProcessPayment` — charges the customer
3. `DispatchToFulfillment` — hands off to the warehouse

**Compensating transactions** (run in reverse if something fails):
- `RefundPayment` — undoes `ProcessPayment`
- `ReleaseInventory` — undoes `ReserveInventory`

There is no compensation for `DispatchToFulfillment` — if dispatch fails, it
never completed, so there's nothing to undo there.

***

### Part B – Call the forward steps

`fulfillment_workflow.go` uses **named return values** — `(result OrderResult, err error)` —
so a plain `return` sends back whatever `result` and `err` currently hold. Call
the three forward steps in order, assigning into `err` and returning
immediately if any of them fails:

```go
if err = workflow.ExecuteActivity(ctx, fa.ReserveInventory, order).Get(ctx, &reservationID); err != nil {
    return
}
if err = workflow.ExecuteActivity(ctx, fa.ProcessPayment, order).Get(ctx, &paymentConfirmation); err != nil {
    return
}
var trackingNumber string
if err = workflow.ExecuteActivity(ctx, fa.DispatchToFulfillment, order, reservationID).Get(ctx, &trackingNumber); err != nil {
    return
}
```

On success, set `result` to a `"FULFILLED"` `OrderResult` (including
`trackingNumber`) and return.

Remove the `_ = fa` / `_ = reservationID` / `_ = paymentConfirmation` placeholder
lines and the `NOT_IMPLEMENTED` result once you're using those variables for real.

***

### Part C – Compensate in reverse order

Go has no `try`/`catch`. A `defer` block right after your state variables are
declared plays the same role: it always runs on the way out of the function,
and — because it's a closure — it can see exactly which steps completed.

```go
defer func() {
    if err == nil {
        return
    }
    if paymentConfirmation != "" {
        _ = workflow.ExecuteActivity(ctx, fa.RefundPayment, paymentConfirmation).Get(ctx, nil)
    }
    if reservationID != "" {
        _ = workflow.ExecuteActivity(ctx, fa.ReleaseInventory, reservationID).Get(ctx, nil)
    }
    result = OrderResult{
        OrderID:             order.OrderID,
        Status:              "FAILED",
        ReservationID:       reservationID,
        PaymentConfirmation: paymentConfirmation,
    }
    err = nil // the saga already handled the failure — return a result, not an error
}()
```

The empty-string guards are the key — if `ReserveInventory` failed before
returning, `reservationID` is still `""`, so there's nothing to release.

> **Go vs. Java:** Java wraps the forward steps in `try` and puts compensation
> in `catch`. Go has no exceptions to catch, so `defer` — which always runs on
> the way out, success or failure — takes on that job instead. Checking
> `err == nil` inside the deferred func is the Go equivalent of "did we get
> here via the catch block."

***

### Part D – Run and observe

1. Click the [button label="Terminal 1 - Worker" background="#444CE7"](tab-1) tab and start the Worker:

   ```bash,run
   go run ./cmd/worker/
   ```

2. Click the [button label="Terminal 2 - Starter" background="#444CE7"](tab-2) tab and run the Starter:

   ```bash,run
   go run ./cmd/starter/
   ```

`DispatchToFulfillment` fails 30% of the time. Run the Starter a few times until
you see a failure. In the [button label="Temporal Web UI" background="#444CE7"](tab-3) tab, open `fulfillment-ORD-5001` and look at the
Event History — you should see `RefundPayment` and `ReleaseInventory` activity
events immediately following the dispatch failure.

**Discussion:** Where does the Saga pattern apply in your own workflows? Think
about: payment + shipping + notification pipelines, multi-service onboarding
flows, reservation + pricing + booking steps.

***

Click **Check** when done, or **Solve** to see the reference solution.
