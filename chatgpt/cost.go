package chatgpt

import (
	"github.com/muzykantov/tgpt/chat"
	"github.com/sashabaranov/go-openai"
)

// CostPer1k represents the cost for 1000 tokens (input & output).
// Tokens are a unit of text used for pricing models where 1,000 tokens
// are approximately equivalent to 750 words. Costs are split into Input
// and Output to differentiate between the text fed into the model and
// the text generated by the model, respectively.
type CostPer1k struct {
	Input  chat.Cost // Cost for input tokens per 1,000 tokens
	Output chat.Cost // Cost for output tokens per 1,000 tokens
}

// Predefined cost structures for various OpenAI models with different
// capabilities and price points. The costs are specified per 1,000 tokens,
// allowing for an easy comparison between the token utilization costs of
// different models. This is especially useful for applications needing
// to budget or make decisions based on the cost efficiency of the models.
var (
	GPT3Dot5TurboCtx4k = CostPer1k{
		Input:  0.0015,
		Output: 0.002,
	} // Cost structure for GPT-3.5 Turbo with a 4k token context.
	GPT3Dot5TurboCtx16k = CostPer1k{
		Input:  0.003,
		Output: 0.004,
	} // Cost structure for GPT-3.5 Turbo with a 16k token context.
	GPT3Dot5Turbo1106Ctx16k = CostPer1k{
		Input:  0.001,
		Output: 0.002,
	} // Cost structure for GPT-3.5 Turbo with a 16k token context.
	GPT4Ctx8k = CostPer1k{
		Input:  0.03,
		Output: 0.06,
	} // Cost structure for GPT-4 with an 8k token context.
	GPT4Ctx32k = CostPer1k{
		Input:  0.06,
		Output: 0.12,
	} // Cost structure for GPT-4 with a 128k token context.
	GPT4Turbo1106Ctx128k = CostPer1k{
		Input:  0.01,
		Output: 0.03,
	} // Cost structure for GPT-4 Turbo with a 128k token context.
)

// Cost provides a mapping from model identifiers to their respective CostPer1k
// structures. This map is used to quickly look up the cost of using a particular
// model based on its identifier, which can be a key factor when selecting
// a model for use in applications.
var Cost = map[string]CostPer1k{
	openai.GPT3Dot5Turbo:    GPT3Dot5TurboCtx4k,      // Maps GPT-3.5 Turbo to its corresponding 4k token context cost structure.
	openai.GPT3Dot5Turbo16K: GPT3Dot5TurboCtx16k,     // Maps GPT-3.5 Turbo with 16k context to its cost structure.
	openai.GPT4:             GPT4Ctx8k,               // Maps GPT-4 with 8k context to its cost structure.
	openai.GPT432K:          GPT4Ctx32k,              // Maps GPT-4 with 32k context to its cost structure.
	"gpt-4-1106-preview":    GPT4Turbo1106Ctx128k,    // Maps GPT-4 Turbo with 128k context to its cost structure.
	"gpt-3.5-turbo-1106":    GPT3Dot5Turbo1106Ctx16k, // Maps GPT-3.5 Turbo with 16k context to its cost structure.
}
