---
slug: parallel-activities
id: vl4lpsmcpxmj
type: challenge
title: 'Exercise 3: Parallel Activities'
teaser: Fan out warehouse checks concurrently using workflow.Go() goroutines and a
  channel.
notes:
- type: text
  contents: |-
    In Exercise 2, the inventory reservation checked warehouses **one at a time** —
    sequential activity calls. At scale (thousands of orders per hour), that is a
    real latency problem.

    **This exercise makes all warehouse checks run simultaneously.**

    In Go SDK, `workflow.Go()` spawns a durable goroutine inside the workflow sandbox.
    `workflow.NewChannel()` lets goroutines pass results back safely.

    Hit **Start** when you're ready.
tabs:
- id: 5jelj1nxcokn
  title: Code
  type: code
  hostname: workshop
  path: /workspace/exercise
- id: fegmmnicvsyt
  title: Terminal 1 - Worker
  type: terminal
  hostname: workshop
  workdir: /workspace/exercise
- id: axvo0xsy6uje
  title: Terminal 2 - Starter
  type: terminal
  hostname: workshop
  workdir: /workspace/exercise
- id: wwlfwt063dys
  title: Temporal Web UI
  type: service
  hostname: workshop
  path: /
  port: 8080
difficulty: basic
timelimit: 3600
enhanced_loading: null
---

## Exercise 3: Parallel Activities

Open **`inventory_workflow.go`** in the [button label="Code" background="#444CE7"](tab-0) tab. Look for `// TODO` comments.

***

### Part A – Create the results channel

Above the loop, create a channel to collect goroutine results:

```go
resultCh := workflow.NewChannel(ctx)
```

***

### Part B – Fan out with `workflow.Go()`

Inside a loop over `Warehouses`, launch one goroutine per warehouse:

```go
for _, warehouseID := range Warehouses {
    wID := warehouseID  // capture loop variable — critical in Go
    workflow.Go(ctx, func(gCtx workflow.Context) {
        var reservationID string
        err := workflow.ExecuteActivity(gCtx, wa.CheckWarehouseInventory, wID, sku, quantity).Get(gCtx, &reservationID)
        if err != nil {
            resultCh.Send(gCtx, "")
            return
        }
        resultCh.Send(gCtx, reservationID)
    })
}
```

> `workflow.Go()` is Temporal's durable goroutine — safe inside workflow code,
> replays correctly. Unlike regular Go goroutines, it participates in the workflow sandbox.

***

### Part C – Collect results

Receive one result per goroutine. Return as soon as you find stock:

```go
for i := 0; i < len(Warehouses); i++ {
    var reservationID string
    resultCh.Receive(ctx, &reservationID)
    if reservationID != "" {
        return reservationID, nil
    }
}
return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")
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

In the [button label="Temporal Web UI" background="#444CE7"](tab-3) tab, open `inventory-ORD-3001`. Notice the `CheckWarehouseInventory`
activity tasks are scheduled with nearly identical timestamps — all running concurrently.
Compare with Exercise 2 where they ran one at a time.

***

Click **Check** when done, or **Solve** to see the reference solution.
