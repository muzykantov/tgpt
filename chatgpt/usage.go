// Package openai provides types and functions to interact with OpenAI's APIs
// for natural language processing and other AI-driven tasks.
package chatgpt

import (
	"fmt"

	"github.com/muzykantov/tgpt/chat"
)

// Usage tracks the token consumption for a single chat request.
// It accounts for both the input and output tokens, thereby providing
// a complete overview of the tokens utilized during an interaction with
// the chat service. This information is crucial for understanding usage
// metrics, estimating costs, and monitoring the efficiency of the token usage.
//
// Input tokens typically correspond to the user's query or prompt, while
// output tokens correspond to the generated response by the assistant.
type Usage struct {
	Input  int // Input represents the number of tokens provided by the user.
	Output int // Output represents the number of tokens in the assistant's response.
}

// CalculateCost calculates the total cost of the tokens consumed during a chat request.
// It considers the number of input and output tokens as detailed in the Usage struct and
// the provided cost per 1,000 tokens, represented by the CostPer1k struct. This function
// is designed to separately compute the costs for input and output tokens and then aggregate
// them to provide a unified cost value represented by the chat.Cost type.
//
// The cost calculation is performed by dividing the total number of tokens (input or output)
// by 1,000 and then multiplying by the corresponding cost per 1,000 tokens (input or output).
// The function returns the sum of the input and output token costs as a single chat.Cost value.
//
// Parameters:
// - costPerToken: A CostPer1k struct containing the cost per 1,000 input and output tokens.
//
// Returns:
// - A chat.Cost value representing the total cost of the consumed tokens.
//
// Note:
// The function assumes the cost per token provided is for 1,000 tokens. If the pricing model
// uses different costs for input and output tokens, this will be accurately reflected in the
// total cost.
func (u *Usage) CalculateCost(costPerToken CostPer1k) chat.Cost {
	inputCost := float64(u.Input) / 1000 * float64(costPerToken.Input)
	outputCost := float64(u.Output) / 1000 * float64(costPerToken.Output)
	return chat.Cost(inputCost + outputCost)
}

// CalculateCostByModel calculates the total cost of the tokens consumed
// during a chat request based on the model name. It uses the global GPTCost
// map to find the CostPer1k structure associated with the given model name and
// then calls CalculateCost method to compute the cost.
//
// modelName: The identifier for the model, which is used to retrieve the cost
// structure from the GPTCost map.
//
// Returns:
// chat.Cost: The total cost for the number of input and output tokens.
// error: An error if the model name does not exist in the GPTCost map or other
// calculation issues occur.
func (u *Usage) CalculateCostByModel(modelName string) (chat.Cost, error) {
	costPer1k, ok := Cost[modelName]
	if !ok {
		return 0, fmt.Errorf("model name '%s' not found in the cost map", modelName)
	}

	return u.CalculateCost(costPer1k), nil
}
