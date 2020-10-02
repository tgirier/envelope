FROM golang:1.14-alpine AS build

WORKDIR /src/cmd
COPY cmd/main.go /src/cmd
COPY *.go go.mod /src/
RUN CGO_ENABLED=0 go build -o /bin/chat

FROM scratch
COPY --from=build /bin/chat /bin/chat
ENTRYPOINT ["/bin/chat"]