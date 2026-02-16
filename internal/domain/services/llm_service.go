package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type LLMService struct {
	client *openai.Client
	model  string
}

func NewLLMService(apiKey, model string) *LLMService {
	client := openai.NewClient(apiKey)
	return &LLMService{
		client: client,
		model:  model,
	}
}

// IntentData represents extracted information from user query
type IntentData struct {
	Entities struct {
		People        []string `json:"people"`
		Organizations []string `json:"organizations"`
		Locations     []string `json:"locations"`
		Products      []string `json:"products"`
	} `json:"entities"`
	Concepts      []string `json:"concepts"`
	Intent        string   `json:"intent"`
	SearchQuery   string   `json:"search_keywords"`
	Category      string   `json:"category,omitempty"`
	Source        string   `json:"source,omitempty"`
	Location      string   `json:"location,omitempty"`
	MinScore      float64  `json:"min_score,omitempty"`
}

// ExtractIntent uses LLM to extract entities, concepts, and intent from user query
func (s *LLMService) ExtractIntent(ctx context.Context, query string, userLat, userLon *float64) (*IntentData, error) {
	systemPrompt := `You are a news query analyzer. Extract ONLY the predefined entities present in the user's query.

ENTITY LABELS TO EXTRACT (real-world objects):
- "people": Names of individuals
- "organizations": Companies, institutions, news sources (e.g., Reuters, BBC, CNN)
- "locations": Geographic places (e.g., Mumbai, Silicon Valley)
- "products": Specific products mentioned

DETERMINE INTENT based on which entities are present:
- "nearby": If query contains a LOCATION entity and user provided their location (user wants news near that place)
- "source": If query contains ORGANIZATION entity that's a news source (user wants articles from that source)
- "category": If query mentions content category keywords like technology, politics, sports, business, world, national
- "score": If query uses importance words: important, trending, top, breaking, urgent, critical
- "search": Default if no specific intent above

CATEGORIES: technology, politics, sports, business, world, national, entertainment, health, science, education

Return ONLY valid JSON in this format:
{
  "entities": {
    "people": ["name1", "name2"],
    "organizations": ["org1", "org2"],
    "locations": ["location1"],
    "products": ["product1"]
  },
  "concepts": ["abstract_idea1", "abstract_idea2"],
  "intent": "nearby|category|source|score|search",
  "search_keywords": "extracted main keywords",
  "category": "category if detected",
  "source": "news source if detected",
  "location": "location if detected for nearby search",
  "min_score": 0.7 if score intent detected
}`

	userLocation := "not provided"
	if userLat != nil && userLon != nil {
		userLocation = fmt.Sprintf("%.4f, %.4f", *userLat, *userLon)
	}

	userPrompt := fmt.Sprintf(`Query: "%s"

User Location: %s (latitude, longitude)

Extract ONLY the predefined entities. Determine intent based on which entities are present.
Return valid JSON only.`, query, userLocation)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
		MaxTokens:   300,
		Temperature: 0.1,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	// Parse JSON response
	var intentData IntentData
	if err := json.Unmarshal([]byte(content), &intentData); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w (content: %s)", err, content)
	}

	return &intentData, nil
}

// GenerateSummary creates a concise summary of an article
func (s *LLMService) GenerateSummary(ctx context.Context, title, description string) (string, error) {
	prompt := fmt.Sprintf(`Summarize this news article in 2-3 concise sentences:

Title: %s
Description: %s

Provide only the summary without any preamble.`, title, description)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		MaxTokens:   150,
		Temperature: 0.3,
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
