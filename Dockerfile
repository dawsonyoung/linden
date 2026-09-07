# syntax=docker/dockerfile:1

# Web assets. src/web has no package.json until Stage C.3; the placeholder keeps
# the final image layout stable so C.3 changes only this stage.
FROM node:22-alpine AS web
WORKDIR /build
COPY src/web/ ./
RUN if [ -f package.json ]; then \
  npm ci && npm run build; \
  else \
  mkdir -p build && echo "web assets arrive in Stage C.3" > build/.placeholder; \
  fi

FROM golang:1.22-alpine AS server
WORKDIR /build
COPY src/go.mod ./
RUN go mod download
COPY src/ ./
COPY --from=web /build/build/ web/build/
ARG VERSION=dev
ARG COMMIT=unknown
# Static binary: the runtime image has no libc and no shell.
RUN CGO_ENABLED=0 GOOS=linux go build \
  -trimpath \
  -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
  -o /out/linden ./cmd/

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=server /out/linden /linden
ENV LINDEN_HOST=0.0.0.0
EXPOSE 8080
USER nonroot:nonroot
# No shell in the image, so the binary probes itself.
HEALTHCHECK --interval=10s --timeout=3s --start-period=3s --retries=3 \
  CMD ["/linden", "-health"]
ENTRYPOINT ["/linden"]
