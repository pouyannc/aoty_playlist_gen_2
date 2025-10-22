FROM golang:1.23.4-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .


FROM alpine:3.20

WORKDIR /root/

RUN apk add --no-cache bind-tools chromium ca-certificates


COPY --from=build /app/server .
COPY dist ./dist
ENV PORT=8080

EXPOSE 8080
CMD ["./server"]