package main

import (
	"context"
	"os"

	elizav1 "connect-examples-go/internal/gen/connectrpc/eliza/v1"

	openai "github.com/sashabaranov/go-openai"
)

// StartOpenAIClient creates a new client goroutine and returns back
// the queue channel for requests.

type OpenAIRequest struct {
	Message *elizav1.ChatRequest
	ClientID string
	ResponseChannel chan OpenAIResponse
	Context context.Context
}

type OpenAIResponse struct {
	Message *elizav1.ChatResponse
}

// StartOpenAIClient creates:
// * a new client goroutine
// * a new openai client goroutine
// and returns back the queue channel for requests.
func StartOpenAIClient() chan OpenAIRequest {
	clientRequests := make(chan OpenAIRequest)
	requests := make(chan OpenAIRequest)
	clientMessages := sync.Map{}
	messageId := 0
	openaiClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	go func() {
		for {

			req, ok := <-clientRequests
			ctx := req.Context
			clientid := GetClientID(ctx)
			count := incCount(clientMessages, clientid, 6)
			if count < 6 {
				req.id = messageId++
				requests <- req
				req.ResponseChannel <- OpenAIResponse{
					Message: &elizav1.ChatResponse{
						Enqueued: &elizav1.ChatEnqueued{
							Status: "Enqueued",
							id: id,
						},
					},
				}
			}
			if count == 6 {
				// send back queue full error
				req.ResponseChannel <- OpenAIResponse{
					Message: &elizav1.ChatResponse{
						Enqueued: &elizav1.ChatEnqueued{
							Status: "Queue is full, please wait",
						},
					},
				}


			if err == nil {
			}
			if !ok {
				// Channel closed, exit goroutine.
				close(requests)
				return
			}
		}
	}()
	// Start a goroutine to receive requests from the client.
	go func() {
		for {
			req, ok := <-requests
			// once the request is processed, remove it from the map
			if !ok {
				// Channel closed, exit goroutine.
				return
			}
			// Send the request to the client.
			clientRequests <- req
		}
	}()
	// Return the channel so callers can queue requests.
	return clientRequests
	}
}
