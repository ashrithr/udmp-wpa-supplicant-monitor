FROM golang:alpine as builder

# Add source code
RUN mkdir /build
ADD . /build/
WORKDIR /build

# Build the code
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /build/main /app/
ENTRYPOINT ["/app/main"]
