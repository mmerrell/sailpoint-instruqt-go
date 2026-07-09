# Exercise 4: Local Activities

In this exercise you will add fast in-process validation steps using local
activities, cutting unnecessary Temporal Actions.

**What you'll learn:**
- How local activities differ from remote activities in execution and history
- When a step belongs in a local activity vs. a regular one
- How to read the Temporal UI to confirm the round-trip was actually skipped

**Time:** ~10 minutes

---

## Background

`ValidateOrder` and `FraudCheck` are pure in-process computations — no network
calls. Running them as regular activities means two unnecessary task queue
round-trips to the Temporal Server per workflow execution. At 20K orders/day
that's 40K wasted Actions. Local activities eliminate those round-trips and
show up in history as a single `MarkerRecorded` event instead of the
schedule/start/complete triplet.

---

## Setup

```bash
temporal server start-dev
```

Open two terminals and `cd` into `4_local_activities/practice/`.

```bash
go mod tidy
```

---

## Part A: Add local activities to `FulfillmentWorkflow`

Open `fulfillment_workflow.go`.

Local activities need their own options type and their own context. Add this
**before** the child workflow call:

```go
lao := workflow.LocalActivityOptions{
    StartToCloseTimeout: 5 * time.Second,
}
localCtx := workflow.WithLocalActivityOptions(ctx, lao)
lfa := &LocalFulfillmentActivities{}

// Validate — returns only error, no result to capture
if err := workflow.ExecuteLocalActivity(localCtx, lfa.ValidateOrder, order).Get(localCtx, nil); err != nil {
    return OrderResult{}, err
}

// Fraud check — returns a risk score string
var riskScore string
if err := workflow.ExecuteLocalActivity(localCtx, lfa.FraudCheck, order).Get(localCtx, &riskScore); err != nil {
    return OrderResult{}, err
}
log.Info("Fraud check passed", "riskScore", riskScore)
```

> **Key difference from remote activities:**
> - `workflow.LocalActivityOptions` not `workflow.ActivityOptions`
> - `workflow.WithLocalActivityOptions()` not `workflow.WithActivityOptions()`
> - `workflow.ExecuteLocalActivity()` not `workflow.ExecuteActivity()`
>
> Everything else — `.Get(ctx, &result)`, error handling — is identical.

---

## Part B: Run it

```bash
# Terminal 1
go run ./cmd/worker/

# Terminal 2
go run ./cmd/starter/
```

---

## What to look for in the UI

Open **http://localhost:8233** and find `fulfillment-ORD-4001`.

| Event | What it means |
|---|---|
| `MarkerRecorded` (×2) | Each local activity — no schedule/start/complete, just one marker |
| `ChildWorkflowExecutionStarted` | Inventory child workflow kicked off |
| `ActivityTaskScheduled/Started/Completed` (×2) | ProcessPayment and DispatchToFulfillment — regular remote activities |

The two local activity calls produced **2 events** total. If they were regular
activities they would have produced **6 events** (3 per activity). At 20K
orders/day, that's 80K fewer Actions.

---

## Solution

The complete solution is in `4_local_activities/solution/`.
