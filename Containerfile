# Build stage
FROM registry.access.redhat.com/ubi9/ubi AS build
RUN dnf install -y go && dnf clean all
WORKDIR /src
COPY . .
RUN go mod init str8edgedave/alert-receiver && \
    go mod tidy && \
    go build .
RUN ls -al /src

# Final image
FROM registry.access.redhat.com/ubi9/ubi-micro
WORKDIR /app
COPY --from=build /src/alert-receiver .
ENV HOST=0.0.0.0
CMD ["./alert-receiver"]

