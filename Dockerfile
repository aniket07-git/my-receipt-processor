# -------------------------------------------------------------------
# 1) Build Stage
#    - We build the Go application in a container with the Go toolchain.
#    - CGO_ENABLED=0 ensures we build a fully static binary (no glibc dependencies).
# -------------------------------------------------------------------
    FROM golang:1.20-bullseye AS builder

    WORKDIR /app
    
    # Copy mod/sum, download deps
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy all source files
    COPY . .
    
    # Disable CGO for a static build; strip symbols to reduce size
    ENV CGO_ENABLED=0
    RUN go build -o /app/receipt-processor -ldflags="-s -w" ./cmd/receipt-processor
    
    # -------------------------------------------------------------------
    # 2) Final Stage
    #    - We use a minimal distroless image that has no shell or package manager.
    #    - Because our binary is fully static, we can also use 'FROM scratch' if preferred.
    # -------------------------------------------------------------------
    FROM gcr.io/distroless/static
    
    # Copy the static binary from the builder stage
    COPY --from=builder /app/receipt-processor /receipt-processor
    
    # Expose port 8080
    EXPOSE 8080
    
    # Run the binary
    ENTRYPOINT ["/receipt-processor"]
    