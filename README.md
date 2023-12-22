Go GRPC Interview
===========

This interview is meant to evaluate someone who is familiar with Golang and GRPC.

This repo contains an example RPC service built with [Connect][connect].
Its API is defined by a [Protocol Buffer schema][schema], and the service
supports the [gRPC][grpc-protocol], [gRPC-Web][grpcweb-protocol], and [Connect
protocols][connect-protocol].

## Installation

**We recommend you set this up before the interview!**

You will need to install the following for this exercise:

```bash
$ go install github.com/bufbuild/buf/cmd/buf@latest
$ go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest # we need this to verify the grpc call
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
```

If you prefer not to install these tools, we can use Replit or CodeSpaces for this interview.

### Check that installation worked

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

# **STOP HERE**! We will go over the rest during the interview.


## Interview Question

Once everything is set up, your task is to implement a new bidirectional stream called `QueueConversation`.
Similar to `Converse`, `QueueConversation` takes in a payload `sentence`, but rather than return a reply right away, queues the sentence.

The queue should have a limit. If the queue is full, the stream should message back that the queue is full, but keep the connection open.


### Software Design and Stubbing out the GRPC calls
In the first 30 minutes, we will design the grpc call and generate the stubs together.

We will then discuss your approach to solving the problem at a high level. By the end of this session we should know the following:
* The rough design of the queue and how you plan on handling the concurrency.
* How do you plan on validating that it works? What checkpoints do you want to make along the way?
* The surface area is quite big (which is on purpose!), so what parts you are interested in and where you intend to focus your time.

### Implementation Portion
After the initial conversation, we will leave you to implement the solution! I will be available via Zoom or Slack (totally your preference) for questions and to help you get unstuck.

The response should look like the following:
```
grpcurl -plaintext -d '@'  localhost:8080 connectrpc.eliza.v1.ElizaService/QueueConversation
{"sentence":"test is a test"}
{
  "sentence": "enqueued"
}
{"sentence":"you gotta"}
{
  "sentence": "enqueued"
}
{"sentence":"you gotta"}
{
  "sentence": "enqueued"
}
{"sentence":"you gotta"}
{"sentence":"you gotta"}
{"sentence":"you gotta"}
{
  "sentence": "queue is full!"
}
{
  "sentence": "queue is full!"
}
{
  "sentence": "queue is full!"
}
```

Another goroutine should dequeue the job and respond back with a call to openai. We will provide you a token. Here is an example call via curl:

```
curl https://api.openai.com/v1/chat/completions   -H "Content-Type: application/json"   -H "Authorization: Bearer $OPENAI_API_KEY"   -d '{
    "model": "gpt-4",
    "messages": [
      {
        "role": "system",
        "content": "You are a poetic assistant, skilled in explaining complex programming concepts with creative flair."
      },
      {
        "role": "user",
        "content": "Compose a poem that explains the concept of recursion in programming."
      }
    ]
  }'
```

This is a sample response:
```
{
  "id": "chatcmpl-8Y6YjZAAIR7RulcCNAJS8kyrzFkcz",
  "object": "chat.completion",
  "created": 1703139057,
  "model": "gpt-4-0613",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "In the realm where logic flows, where codes create and nerds compose,\nThere lies a spell both sleek and clean, like a two-way mirror in machine,\nIn its essence pure and serenely mien, we behold the concept named Recursion.\n\nA function calling itself again, a loop unbroken, no end to attain,\nIt dives within its own domain, a reflection questing its own version.\nSelf-contained, a paradox in motion, an infinity mirror's reflection,\nThis dear is the charm, a jest in the name of Recursion.\n\nThe base case is the boundary end, the point at which the path will bend,\nFor the cycle to prevent the never-ending grind, Solutions, you see, it must always find.\nA condition met to stop the spin, a promise to the chaos, that order will win,\nIn these simple rules and not a single diversion, lies the elegant beauty called Recursion.\n\nFactorials, Fibonacci, and more, the traces of this method vast to explore,\nDizzying towers of Hanoi's game, trees of binary know its fame,\nA sequence by nature, or by the sorting game, recursion assists them just the same.\nIn this vast expanse of logical immersion, lives the dynamic power of Recursion.\n\nRecursion, the serpent eating its tail, going deeper without fail,\nTraveling down the rabbit hole, finding its base to reach the goal,\nReturns upon itself, a circle, not a line, a puzzle, a riddle, a sign,\nA tool of the code crafted in such precision, Oh marvel at the thought of Recursion.\n\nTo understand, to learn, to know this art, begin at the end, that’s where it starts,\nThe smallest fragment holds the key, to unlock the secrets for you and me,\nArmed with this knowledge, with this vision, we unmask the magic of Recursion.\n\nAn epitome of complexity and reduction, is indeed our friend, sweet Recursion,\nIn layers of code, it weaves its verse, an ode to order from the universe.\nSo toast to the code that spirals in precision, and raises a hat to lovely Recursion."
      },
      "logprobs": null,
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 39,
    "completion_tokens": 436,
    "total_tokens": 475
  },
  "system_fingerprint": null
}
```

You can find the documentation here: https://platform.openai.com

After that, make a pool of gorountines that is configurable and dequeues.

## Extra Time
If there is extra time, it would be good to add retries and ratelimiting/circuit breaking between the queue and the call to OpenAI, as OpenAI can be very slow.

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
