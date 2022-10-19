FROM golang:1.19-alpine as build

WORKDIR $GOPATH/app/

RUN apk add git

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/

#resulting app
FROM scratch as final
COPY --from=build go/app/main /app/
WORKDIR /app
ENTRYPOINT [ "./main" ]