# --- build stage ---
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /notiflex-api .

# --- runtime stage ---
FROM scratch
COPY --from=build /notiflex-api /notiflex-api
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/notiflex-api"]
