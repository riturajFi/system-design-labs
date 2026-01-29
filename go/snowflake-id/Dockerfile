# syntax=docker/dockerfile:1
# Tells Docker which Dockerfile “language”/feature set to use.
# `docker/dockerfile:1` enables modern Dockerfile features and predictable parsing.

FROM golang:1.22 AS build
# `FROM` picks the base image for this stage.
# `golang:1.22` is an image that already has Go installed so we can compile the program.
# `AS build` names this stage "build" so we can copy artifacts from it later (multi-stage build).

WORKDIR /src
# Sets the working directory inside the container for subsequent commands.
# If /src doesn’t exist, Docker creates it.
# After this, `COPY` and `RUN` paths are relative to /src unless you use absolute paths.

COPY go.mod ./
# Copies the local file `go.mod` from your project folder (build context) into the container at `/src/go.mod`.
# Doing this first helps Docker cache dependency/module steps (in larger projects you’d also COPY go.sum).

COPY cmd ./cmd
# Copies your local `cmd` directory into the container at `/src/cmd`.
# This includes `cmd/idgen` which contains the Go "main" package we build.

COPY internal ./internal
# Copies your local `internal` directory into the container at `/src/internal`.
# This includes the library code imported by `cmd/idgen`.

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/idgen ./cmd/idgen
# `RUN` executes a command *during image build time* and bakes the results into the image layer.
# `CGO_ENABLED=0` disables C bindings so the binary can be fully static-friendly and simpler to run.
# `go build` compiles the program.
# `-trimpath` removes local filesystem paths from the compiled binary for cleaner/reproducible builds.
# `-ldflags="-s -w"` strips symbol/debug info to make the binary smaller.
# `-o /out/idgen` writes the compiled binary to `/out/idgen` inside the build stage.
# `./cmd/idgen` is the package path to build (your server program).

FROM gcr.io/distroless/static:nonroot
# Starts a *new stage* for the final runtime image.
# `distroless/static:nonroot` is a minimal image intended only to run a binary:
# - no package manager, no shell, far fewer files (smaller + fewer attack surface)
# - `nonroot` means it runs as a non-root user by default (safer).

COPY --from=build /out/idgen /idgen
# Copies a file from the earlier stage named `build` into this final image.
# `--from=build` means “take it from the build stage filesystem”.
# `/out/idgen` is the compiled binary we produced.
# `/idgen` is where we place it in the final image.

ENTRYPOINT ["/idgen"]
# Defines what runs when the container starts.
# JSON array form means:
# - Docker runs it directly (no shell involved)
# - avoids quoting issues and is the recommended form
# So: `docker run ... snowflake-id` will execute `/idgen` by default.
