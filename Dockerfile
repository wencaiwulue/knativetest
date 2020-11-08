FROM golang:latest

WORKDIR app/

COPY go.* ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . ./
RUN CG0_ENABLED=0 GOOS=linux go build -mod=readonly -v -o knative test/knative/

ENTRYPOINT ["./knative/knative"]