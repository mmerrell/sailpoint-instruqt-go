# Temporal Go Hands-On Lab

Instruqt track: five progressive coding challenges using the Temporal Go SDK.
Built on a durable order-fulfillment scenario (inventory reservation, payment, dispatch).

## Track Challenges

| # | Slug | Concept | Key APIs |
|---|------|---------|----------|
| 1 | `converting` | Replace retry loops with Activities + Workflow | `ExecuteActivity`, `ActivityOptions`, `ApplicationError` |
| 2 | `child-workflows` | Decompose into child workflow with separate history | `ExecuteChildWorkflow`, `ChildWorkflowOptions` |
| 3 | `parallel-activities` | Fan out warehouse checks concurrently | `workflow.Go()`, `workflow.Channel` |
| 4 | `local-activities` | Move fast steps to Local Activities | `ExecuteLocalActivity`, `LocalActivityOptions` |
| 5 | `saga` | Compensating transactions for multi-step failure recovery | reverse-order compensation, null-guarded state |

## Prerequisites (self-paced, without Instruqt)

- Go 1.22+
- Temporal CLI: `brew install temporal` (or download from temporal.io)
- A running local Temporal server: `temporal server start-dev`

## Repo Layout

```
temporal-go-hands-on/
├── docker/
│   └── Dockerfile                  # golang:1.23 + Temporal CLI, exercises baked in
├── 1_converting/                   # Exercise source, baked into the sandbox image
│   ├── practice/                   # Learner starting point (// TODO stubs)
│   └── solution/                   # Reference implementation
├── 2_child_workflows/
├── 3_parallel_activities/
├── 4_local_activities/
├── 5_saga/
├── track/                          # Instruqt track definition
│   ├── track.yml
│   ├── config.yml                  # Container config
│   ├── track_scripts/
│   │   └── setup-workshop          # Starts the Temporal dev server
│   ├── 01-converting/
│   │   ├── assignment.md           # Learner instructions + tab definitions
│   │   ├── config.yml
│   │   ├── setup-workshop
│   │   ├── check-workshop
│   │   ├── solve-workshop
│   │   └── cleanup-workshop
│   ├── 02-child-workflows/
│   ├── 03-parallel-activities/
│   ├── 04-local-activities/
│   └── 05-saga/
└── .github/workflows/build-image.yml  # Rebuilds sandbox image on push to main
```

Each exercise directory has a `practice/` (with TODO comments), a `solution/`,
and a `README.md` for working through it outside of Instruqt.

## The Domain

An order fulfillment pipeline:

1. **Validate + fraud check** — fast in-process checks (local activities in Exercise 4)
2. **Reserve inventory** — check warehouses for stock (child workflow; parallel in Exercise 3)
3. **Process payment** — charge the customer
4. **Dispatch** — hand off to the fulfillment center
5. **Compensate on failure** — undo payment and inventory steps that already completed (Exercise 5)

## Container Architecture

Each challenge runs in a single **container** sandbox:

| Container | Image | Role |
|-----------|-------|------|
| `workshop` | `ghcr.io/mmerrell/temporal-go-sandbox:latest` | Go 1.23, Temporal CLI, exercise code, dev server |

The Temporal dev server is started inside the container at track start
(`temporal server start-dev --ui-port 8080`). Learners get four tabs per
challenge: a native code editor, a terminal for the Worker, a terminal for
the Starter, and a Temporal Web UI service tab.

## Building the Image Locally

```bash
# From the repo root
docker build -f docker/Dockerfile -t temporal-go-sandbox:local .

# Verify the contents
docker run -it --rm temporal-go-sandbox:local bash
# Inside: go version, temporal --version, ls /opt/exercises/
```

## CI/CD

On every push to `main` that touches `docker/Dockerfile` or any exercise
directory, GitHub Actions rebuilds and pushes:
- `ghcr.io/mmerrell/temporal-go-sandbox:latest`
- `ghcr.io/mmerrell/temporal-go-sandbox:<sha>` (pinned reference)

## Pushing to Instruqt

```bash
cd track/
instruqt track push
```

**First push:** the `id:` fields in `track.yml` and each `assignment.md` are
populated by Instruqt. Pull the track after first push to capture the
generated IDs:

```bash
instruqt track pull
```

## Check Script Strategy

All check scripts use **source grep + `go build`** — no end-to-end workflow
execution. Checks are fast and don't flake on timing. The tradeoff is that a
learner could pass checks with syntactically correct but logically wrong
code; the grep checks are surgical enough to catch the common failure modes.
