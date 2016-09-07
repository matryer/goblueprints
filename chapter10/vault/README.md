
## Usage

Run the vaultd service:

```
cd vault/cmd/vaultd
go run main.go
```

Hash a password:

```bash
curl -XPOST -d'{"password":"MySecretPassword123"}' localhost:8080/hash
```

```json
{"hash":"$2a$10$L/Riz9xbgTBDn7F6uLInq.9Tr67PvBCmxzrLgemitnRM53ht7LGpC"}
```

Validate passwords with hashes:

```bash
curl -XPOST -d'{"password":"MySecretPassword123","hash":"$2a$10$L/Riz9xbgTBDn7F6uLInq.9Tr67PvBCmxzrLgemitnRM53ht7LGpC"}' localhost:8080/validate
```

```json
{"valid":true}
```

or if you get the password wrong:

```bash
curl -XPOST -d'{"password":"NOPE","hash":"$2a$10$L/Riz9xbgTBDn7F6uLInq.9Tr67PvBCmxzrLgemitnRM53ht7LGpC"}' localhost:8080/validate
```

```json
{"valid":false}
```

### Compiling protobuf

Install proto3 from source:

```
brew install autoconf automake libtool
git clone https://github.com/google/protobuf
./autogen.sh ; ./configure ; make ; make install
```

Update protoc Go bindings:

```
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

See also https://github.com/grpc/grpc-go/tree/master/examples

Compile the protobuf (from inside `pb` folder):

```
protoc vault.proto --go_out=plugins=grpc:.
```
