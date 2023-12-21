Go GRPC Interview
===========

This interview is meant to evaluate someone who is familiar with Golang and GRPC.

This repo contains an example RPC service built with [Connect][connect].
Its API is defined by a [Protocol Buffer schema][schema], and the service
supports the [gRPC][grpc-protocol], [gRPC-Web][grpcweb-protocol], and [Connect
protocols][connect-protocol].

## Installation

You will need to install the following:

```bash
$ go install github.com/bufbuild/buf/cmd/buf@latest
$ go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest # we need this to verify the grpc call
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
```

## Check that installation work

Run the server like so:
```
go run cmd/demoserver/main.go
```

In another terminal, verify the server is working by:
```bash
curl --header "Content-Type: application/json" \
    --data '{"sentence": "I feel happy."}' \
    localhost:8080/connectrpc.eliza.v1.ElizaService/Sayd
```

Verify grpc streaming is working using [`grpcurl`][grpcurl] and the gRPC protocol:

```bash
grpcurl -plaintext -d '@'  localhost:8080 connectrpc.eliza.v1.ElizaService/Converse
```

Send the following requests one at a time
```
{"sentence":"test"}
{"sentence":"test"}
{"sentence":"bye"}
{"sentence":"bye"}
```

You should see the following output:
```bash
➜  connectrpc-example-go git:(main) ✗ grpcurl -plaintext -d '@'  localhost:8080 connectrpc.eliza.v1.ElizaService/Converse
{"sentence":"test"}
{
  "sentence": "I see. And what does that tell you?"
}
{"sentence":"test"}
{
  "sentence": "Can you elaborate on that?"
}
{"sentence":"bye"}
{
  "sentence": "Goodbye. I'm looking forward to our next session."
}


{"sentence":"bye"}

```

## Legal

Offered under the [Apache 2 license][license].

[blog]: https://buf.build/blog/connect-a-better-grpc
[connect]: https://github.com/connectrpc/connect-go
[connect-protocol]: https://connectrpc.com/docs/protocol
[docs]: https://connectrpc.com
[eliza]: https://en.wikipedia.org/wiki/ELIZA
[grpc-protocol]: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
[grpcurl]: https://github.com/fullstorydev/grpcurl
[grpcweb-protocol]: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md
[license]: https://github.com/connectrpc/examples-go/blob/main/LICENSE.txt
[schema]: https://github.com/connectrpc/examples-go/blob/main/proto/connectrpc/eliza/v1/eliza.proto
