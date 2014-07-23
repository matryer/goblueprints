
## Installing

This section covers what is required to install dependencies for this project.

### Install nsq

```
brew install nsq
```

  * See [http://nsq.io/deployment/installing.html](http://nsq.io/deployment/installing.html)

### Get `go-nsq`

```
go get github.com/bitly/go-nsq
```

## Starting

  # Start `nsqlookupd`
  # Start `nsqd --lookupd-tcp-address=127.0.0.1:4160`
