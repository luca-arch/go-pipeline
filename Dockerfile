#########################################
# Step 1: Modules caching               #
#########################################
FROM golang:1.18-alpine as modules

COPY go.mod go.sum /modules/

WORKDIR /modules
RUN go mod download

#########################################
# Step 2: Builder                       #
#########################################
FROM golang:1.18-alpine as builder

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=0

COPY --from=modules /go/pkg /go/pkg
COPY . /build

WORKDIR /build
RUN go build -o /bin/pipeline ./cmd

#########################################
# Step 3: Final                         #
#########################################
FROM busybox

COPY --from=builder /bin/pipeline /bin/pipeline
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/bin/pipeline"]
