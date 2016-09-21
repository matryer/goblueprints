FROM golang
ADD . /go/src/github.com/matryer/goblueprints/chapter10/vault
RUN cd /go/src/github.com/matryer/goblueprints/chapter10/vault/cmd/vaultd ; go get 
RUN go install github.com/matryer/goblueprints/chapter10/vault/cmd/vaultd
ENTRYPOINT /go/bin/vaultd
EXPOSE 8080 8081