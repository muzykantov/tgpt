package chatgpt

// RequestParams defines the set of parameters used to customize
// an OpenAI request. These parameters allow for tuning the
// behavior of the model during the conversation.
type RequestParams struct {
	MaxTokens        int     // MaxTokens is the maximum number of tokens to generate in a response.
	Temperature      float32 // Temperature adjusts the randomness of the model's output.
	TopP             float32 // TopP is the sampling value to influence token choice; lower values make output more deterministic.
	PresencePenalty  float32 // PresencePenalty adjusts the model to prefer tokens from the input.
	FrequencyPenalty float32 // FrequencyPenalty adjusts the model to avoid tokens from the input.
}

// DefaultRequestParams is a predefined set of parameters representing default
// values for an OpenAI request. It can be used as a starting point when
// configuring a Chat instance.
var DefaultRequestParams = RequestParams{
	MaxTokens:        0,
	Temperature:      0.0,
	TopP:             0.0,
	PresencePenalty:  0.0,
	FrequencyPenalty: 0.0,
}
