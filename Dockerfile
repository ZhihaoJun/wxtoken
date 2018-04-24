FROM zhihaojun/bkgolang
MAINTAINER January

COPY src/ /app/src/
RUN go build -o /main src/app/main.go
ENTRYPOINT ["/main"]
