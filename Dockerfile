FROM golang:alpine AS build_base

RUN apk add --no-cache git && \
    apk add openssl

WORKDIR /tmp/jobsity-code-challenge

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

#RUN ["./certificate.sh"]

RUN CGO_ENABLED=0 go test -v


RUN go build -o ./out/app .

FROM alpine:latest

ENV GIN_MODE=release

COPY --from=build_base /tmp/jobsity-code-challenge/out/app /app/jobsity-code-challenge
#COPY --from=build_base /tmp/jobsity-code-challenge/localhost.crt /app
#COPY --from=build_base /tmp/jobsity-code-challenge/localhost.key /app

COPY public /app/public
COPY config.yml /app

WORKDIR /app

ENV DB_PASSWORD postgres
ENV SECRET_KEY 1wSa36eEa7t7loL7HztSujZKa9n0M7EVyrv6YJAZDt9I5FuvLvSNVxRI5CUdbd0

#EXPOSE 443
EXPOSE 80

ENTRYPOINT ["/app/jobsity-code-challenge"]