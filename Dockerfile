FROM docker:latest

WORKDIR app/

COPY knative/knative ./knative

#RUN go env -w GOPROXY=https://goproxy.cn,direct
#RUN go mod download

#COPY . ./
#RUN CG0_ENABLED=0 GOOS=linux go build -mod=readonly -v -o knative test/knative/
#RUN CG0_ENABLED=0 GOOS=linux go build -v -o knative ./

ENTRYPOINT ["./knative"]