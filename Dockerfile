FROM golang:1.16-alpine

COPY . .
# RUN rm -f go.mod
ENV GOPATH=
# RUN go mod download
#RUN go mod init test
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build -o signedurl

CMD ["/signedurl"]