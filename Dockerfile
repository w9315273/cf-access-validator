FROM golang:1.21-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN go build -trimpath -ldflags "-s -w" -o /out/cf-auth ./apps/cf-auth

FROM gcr.io/distroless/static:nonroot

ENV ADDR="0.0.0.0:9000" \
    TEAM_DOMAIN="<Team>.cloudflareaccess.com" \
    APP_MAP="test=1234567890"

COPY --from=build /out/cf-auth /cf-auth

USER nonroot:nonroot
ENTRYPOINT ["/cf-auth"]