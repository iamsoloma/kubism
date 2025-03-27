FROM golang:1.24.1-alpine AS builder

WORKDIR /build

#Dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

#Build
COPY . .
RUN go build -o ./bin/kubism



FROM alpine:3.21.3 as runner

WORKDIR /workspace

#Copy
COPY --from=builder /build/bin/kubism ./kubism


#Run
EXPOSE 1304-1304
CMD ["./kubism"]