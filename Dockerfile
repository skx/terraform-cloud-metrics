#
# Trivial Dockerfile for this application.
#
# Builds the application, and can be used to expose the HTTP server on
# port 8090 afterward.
#


# STEP1 - Build-image
###########################################################################
FROM golang:alpine AS builder

LABEL org.opencontainers.image.source=https://github.com/metacore-games/terraform-cloud-metrics

# Ensure we have git
RUN apk update && apk add --no-cache git

# Create a working-directory
WORKDIR $GOPATH/src/github.com/metacore-games/terraform-cloud-metrics

# Copy the source to it
COPY . .

# Build the binary.
RUN go build -o /go/bin/terraform-cloud-metrics



# STEP2 - Deploy-image
###########################################################################
FROM alpine

# Create a working directory
WORKDIR /app

# Copy the binary.
COPY --from=builder /go/bin/terraform-cloud-metrics /app/

# Create a group and user
RUN addgroup app && adduser -D -G app -h /app app

# Ensure we run as that non-root user
USER app

# Set CWD
WORKDIR /app

# Expose the port
EXPOSE 8090

# Entrypoint
ENTRYPOINT ["/app/terraform-cloud-metrics"]
