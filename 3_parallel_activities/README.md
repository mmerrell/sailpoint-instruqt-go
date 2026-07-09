# Exercise 3: Parallel Activities

In this exercise you will make the inventory reservation fan out to all
warehouses simultaneously instead of checking them one at a time.

**What you'll learn:**
- How to run activities in parallel using `workflow.Go()` and `workflow.Channel`
- Why sequential external calls become a latency problem at scale
- How to read the Temporal UI to verify the fan-out actually happened

**Time:** ~10 minutes

---

## Background

In Exercise 2, the inventory reservation workflow checked warehouses **one at a
time** — sequential activity calls. At scale (thousands of orders per hour),
checking each warehouse sequentially is a real latency problem. The fix: fan
out in parallel.

---

## Setup

```bash
temporal server start-dev
```

Open two terminals and `cd` into `3_parallel_activities/practice/`.

```bash
go mod tidy
```

---

## Part A: Fan out with `workflow.Go()`

Open `inventory_workflow.go`.

In the Go SDK, `workflow.Go()` spawns a durable goroutine inside the workflow
sandbox. It is **not** a regular Go goroutine — it replays correctly and is safe
to use in workflow code. Use `workflow.NewChannel()` to pass results back to the
main goroutine.

First, create the results channel above your loop:

```go
resultCh := workflow.NewChannel(ctx)
```

Then, inside a loop over `Warehouses`, launch one goroutine per warehouse:

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

> **Go vs. Java:** Java SDK uses `Async.function()` to opt in to concurrency
> (synchronous by default). Go SDK uses `workflow.Go()` goroutines. The mental
> model is the same: fire multiple futures, collect results.

---

## Part B: Collect results

After the launch loop, receive one result per goroutine. Return as soon as you
find stock — don't wait for all warehouses to finish:

```go
for i := 0; i < len(Warehouses); i++ {
    var reservationID string
    resultCh.Receive(ctx, &reservationID)
    if reservationID != "" {
        return reservationID, nil
    }
}
```

> **Go vs. Java:** Java uses `Promise.allOf(promises).get()` to wait for all,
> then iterates. Go's `channel.Receive()` blocks until the next result arrives —
> you can return early as soon as you get a hit without waiting for the rest.

---

## Part C: Handle the out-of-stock case

If all goroutines returned empty strings, no warehouse had stock:

```go
return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")
```

---

## Part D: Run it

```bash
# Terminal 1
go run ./cmd/worker/

# Terminal 2
go run ./cmd/starter/
```

---

## What to look for in the UI

Open **http://localhost:8233** and find `inventory-ORD-3001`.

In the event history, look at the `ActivityTaskScheduled` events for
`CheckWarehouseInventory`. In Exercise 2 they appeared one after another (each
waited for the previous). Here they should appear with nearly identical
timestamps — all scheduled within the same workflow task.

| Exercise 2 (sequential) | Exercise 3 (parallel) |
|---|---|
| Warehouse 1 scheduled → completed → Warehouse 2 scheduled… | All warehouses scheduled in the same workflow task |
| Total time ≈ N × single-warehouse latency | Total time ≈ single-warehouse latency |

---

## Solution

The complete solution is in `3_parallel_activities/solution/`. Diff
`solution/inventory_workflow.go` against Exercise 2's version to see the
parallel channel pattern side by side with the sequential one.
