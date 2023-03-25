FROM golang:buster as dev

WORKDIR /usr/src/app
RUN go install github.com/codegangsta/gin@latest
COPY go.mod go.sum /usr/src/app/
RUN go mod tidy
COPY . . 
RUN go build -o /usr/local/bin/kitchen

FROM busybox

COPY --from=dev /usr/local/bin/kitchen /usr/local/bin/kitchen

EXPOSE 8080
CMD ["/usr/local/bin/kitchnen"]