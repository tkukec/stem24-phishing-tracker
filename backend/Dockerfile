FROM golang:alpine3.18

RUN apk update && apk upgrade && apk add --no-cache --update \
    openssh-client \
    bash \
    git \
    ca-certificates

WORKDIR /app

COPY . .


RUN mkdir -p -m 0600 /root/.ssh && echo "StrictHostKeyChecking no" >> /root/.ssh/config
RUN --mount=type=ssh GO111MODULE=on go mod download -x

# Add default user
RUN adduser -u 1001 -S -G root -h /tmp -s /sbin/nologin \
        -D -H -g "Default Application User" default

EXPOSE 8080

# In order to drop the root user, we have to make some directories world
# writable as OpenShift default security model is to run the container under
# random UID.
RUN chown -R 1001:0 /app

ENTRYPOINT go run cmd/stem24-backend/main.go

# Run container by default as user with id 1001 (default)
USER 1001