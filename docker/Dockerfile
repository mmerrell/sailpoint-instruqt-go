# temporal-go-sandbox
# Base image for the Temporal Go SDK hands-on lab (Instruqt).
#
# Contains:
#   - Go 1.23 (bookworm base)
#   - Temporal CLI (latest)
#   - All five exercise directories baked into /opt/exercises/
#   - Go module cache pre-warmed for all five exercises
#
# The Temporal server is started by the Instruqt track setup script using
# `nohup temporal server start-dev`. This image is for the learner's sandbox only.

FROM golang:1.23-bookworm

# ─── System packages ──────────────────────────────────────────────────────────
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    wget \
    git \
    vim \
    nano \
    less \
    unzip \
    ca-certificates \
    bash-completion \
    && rm -rf /var/lib/apt/lists/*

# ─── Temporal CLI ─────────────────────────────────────────────────────────────
# Pinned release tarball, not the curl|sh installer — the installer's install
# path/shell assumptions aren't reliable inside a non-interactive Docker build.
ARG TEMPORAL_CLI_VERSION=1.3.0
RUN curl -sSL \
    "https://github.com/temporalio/cli/releases/download/v${TEMPORAL_CLI_VERSION}/temporal_cli_${TEMPORAL_CLI_VERSION}_linux_amd64.tar.gz" \
    | tar -xz -C /usr/local/bin temporal \
    && chmod +x /usr/local/bin/temporal \
    && temporal --version

# ─── Non-root user ────────────────────────────────────────────────────────────
RUN useradd -m -s /bin/bash -u 1001 user
WORKDIR /home/user

# ─── Bake in exercise source code ─────────────────────────────────────────────
# The exercise directories at the repo root must contain:
#   1_converting/practice/ and 1_converting/solution/
#   2_child_workflows/practice/ and 2_child_workflows/solution/
#   3_parallel_activities/practice/ and 3_parallel_activities/solution/
#   4_local_activities/practice/ and 4_local_activities/solution/
#   5_saga/practice/ and 5_saga/solution/
COPY 1_converting/ /opt/exercises/1_converting/
COPY 2_child_workflows/ /opt/exercises/2_child_workflows/
COPY 3_parallel_activities/ /opt/exercises/3_parallel_activities/
COPY 4_local_activities/ /opt/exercises/4_local_activities/
COPY 5_saga/ /opt/exercises/5_saga/
RUN chown -R user:user /opt/exercises

# ─── Pre-warm Go module cache ─────────────────────────────────────────────────
# Downloads all dependencies for all ten go.mod files (practice + solution)
# so learner containers start without downloading anything.
USER user
RUN for dir in \
      /opt/exercises/1_converting/practice \
      /opt/exercises/1_converting/solution \
      /opt/exercises/2_child_workflows/practice \
      /opt/exercises/2_child_workflows/solution \
      /opt/exercises/3_parallel_activities/practice \
      /opt/exercises/3_parallel_activities/solution \
      /opt/exercises/4_local_activities/practice \
      /opt/exercises/4_local_activities/solution \
      /opt/exercises/5_saga/practice \
      /opt/exercises/5_saga/solution; \
    do \
      echo "Warming Go module cache for $dir..." && \
      (cd "$dir" && go mod download); \
    done

# ─── Shell environment ────────────────────────────────────────────────────────
RUN echo 'export PATH=$PATH:/usr/local/go/bin:/home/user/go/bin' >> /home/user/.bashrc \
    && echo 'export TEMPORAL_ADDRESS="${TEMPORAL_ADDRESS:-127.0.0.1:7233}"' >> /home/user/.bashrc \
    && echo 'alias tl="temporal workflow list --address ${TEMPORAL_ADDRESS}"' >> /home/user/.bashrc

USER root

CMD ["/bin/bash"]
