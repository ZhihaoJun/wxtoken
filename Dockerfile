FROM zhihaojun/bkgolang

COPY src/ /app/src/
RUN go build -o /main src/app/main.go
ENTRYPOINT ["/main"]
