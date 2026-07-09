# Exercise 2: Child Workflows

In this lab you will extract inventory reservation into a dedicated child workflow.

**What you'll learn:**
- When to use a child workflow vs. a plain activity
- How to pass arguments to a child workflow and collect its result
- How child workflows appear in the Temporal UI with their own separate event history

**Time:** ~10 minutes

---

## Background

The starting code calls `CheckWarehouseInventory` as a plain activity directly in
`FulfillmentWorkflow`. This works, but it means inventory reservation is just a
step in the parent's history — no isolation, no independent retry boundary, no
separate workflow ID you can query.

A child workflow fixes that. The inventory reservation gets:
- Its own workflow ID (`inventory-ORD-xxxx`) — queryable and cancellable independently
- Its own event history — keeps the parent history clean
- Its own retry boundary — a bug in reservation logic doesn't pollute the parent

---

## Setup

Start a local Temporal server if you don't already have one running:

```bash
temporal server start-dev
```

Open two terminals and `cd` into `2_child_workflows/practice/`.

Install dependencies (first time only):

```bash
go mod tidy
```

---

## Part A: Implement `InventoryReservationWorkflow`

Open `inventory_workflow.go`.

The TODO asks you to iterate over `Warehouses` and call `CheckWarehouseInventory`
for each warehouse until you find one with stock.

```go
for _, warehouseID := range Warehouses {
    var reservationID string
    err := workflow.ExecuteActivity(ctx, wa.CheckWarehouseInventory, warehouseID, sku, quantity).Get(ctx, &reservationID)
    if err != nil {
        return "", err
    }
    if reservationID != "" {
        return reservationID, nil
    }
}
```

If no warehouse has stock, return a non-retryable error:

```go
return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")
```

---

## Part B: Call the child workflow from `FulfillmentWorkflow`

Open `fulfillment_workflow.go`.

Replace the `// REMOVE THIS BLOCK` section with a child workflow call:

```go
childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
    WorkflowID: "inventory-" + order.OrderID,
})
var reservationID string
if err := workflow.ExecuteChildWorkflow(childCtx, InventoryReservationWorkflow,
    order.ItemSKU, order.Quantity).Get(ctx, &reservationID); err != nil {
    return OrderResult{}, err
}
```

> **Go vs. Python:** Python uses `await workflow.execute_child_workflow(ChildWorkflow, args=[...])`.
> Go uses `workflow.ExecuteChildWorkflow(...).Get(ctx, &result)` — the `.Get()` is the await.

---

## Part C: Run it

```bash
# Terminal 1 — start the worker
go run ./cmd/worker/

# Terminal 2 — start a workflow
go run ./cmd/starter/
```

Open the Temporal UI at **http://localhost:8233**.

Find the `fulfillment-ORD-1001` workflow. Then find `inventory-ORD-1001` — notice it
has its **own separate event history**, completely independent of the parent.

---

## What to look for in the UI

| Parent workflow | Child workflow |
|---|---|
| `fulfillment-ORD-1001` | `inventory-ORD-1001` |
| Shows `ChildWorkflowExecutionStarted` | Shows individual `ActivityTaskScheduled` events per warehouse |
| Waits on `ChildWorkflowExecutionCompleted` | Has its own Start → Complete timeline |

---

## Solution

The complete solution is in `2_child_workflows/solution/`. Compare your implementation
if you get stuck, or run the solution directly to see the expected behavior before you
start.
