package usecase

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/chanombude/twitter-go-api/internal/db"
	"github.com/pgvector/pgvector-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/genai"
)

func (u *AssistantUsecase) Chat(ctx context.Context, input AssistantInput) (io.Reader, error) {
	if u.config.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  u.config.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	var timelineContext string
	if !u.config.EnableRAG {
		log.Debug().Msg("RAG is explicitly disabled in config")
		timelineContext = "Timeline context is currently disabled."
	} else {
		// 1. Vectorize the query
		var dim int32 = 768

		queryContent := []*genai.Content{
			{Parts: []*genai.Part{{Text: input.Query}}},
		}

		embedRes, err := client.Models.EmbedContent(ctx, u.config.GeminiEmbeddingModel, queryContent, &genai.EmbedContentConfig{
			TaskType:             "RETRIEVAL_QUERY",
			OutputDimensionality: &dim,
		})

		if err != nil || len(embedRes.Embeddings) == 0 {
			log.Error().Err(err).Msg("failed to vectorize query for RAG, falling back to base chat")
			timelineContext = "Timeline context unavailable (vectorization failed)."
		} else {
			queryVector := pgvector.NewVector(embedRes.Embeddings[0].Values)

			// 2. Retrieve Context (RAG)
			rows, err := u.store.ListRelatedTweetsByEmbedding(ctx, db.ListRelatedTweetsByEmbeddingParams{
				Column1: &queryVector,
				Limit:   5,
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to retrieve context from DB for RAG")
				timelineContext = "Timeline context unavailable (DB error)."
			} else if len(rows) == 0 {
				timelineContext = "No recent matching posts found in the timeline."
			} else {
				var contextParts []string
				for _, row := range rows {
					contextParts = append(contextParts, fmt.Sprintf("- %s", row.Content))
				}
				timelineContext = strings.Join(contextParts, "\n")
			}
		}
	}

	// 3. Construct System Prompt
	systemPrompt := `You are the helpful, professional, and polite AI assistant for "Chanom Twitter", a social media platform and Twitter clone. 
    Your goal is to provide accurate, concise, and friendly information to users of Chanom Twitter. 
    You have access to the user's recent [Chat History] and [Timeline Context] from their feed. 
    Prioritize the Timeline Context when answering questions about the app's content or the user's recent experience. 
    Maintain conversational continuity and a helpful tone consistent with a premium social experience.
    CRITICAL: Your memory is limited to the provided history. If asked about older context you don't have, politely state you don't have access to that part of the conversation and ask for clarification.
	` + "\n\n[Timeline Context]\n" + timelineContext

	// 4. Prepare Chat History
	var contents []*genai.Content
	start := 0
	if len(input.History) > 30 {
		start = len(input.History) - 30
	}

	for i := start; i < len(input.History); i++ {
		h := input.History[i]
		role := h.Role
		if role == "assistant" || role == "model" {
			role = "model"
		} else {
			role = "user"
		}
		contents = append(contents, &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: h.Text}},
		})
	}

	// 5. Append the latest query to the chat history
	contents = append(contents, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: input.Query}},
	})

	// 6. Set System Prompt for the model (using GenerateContentConfig)
	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemPrompt}}},
	}

	// 7. Create Pipe for streaming
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		// 8. Use Go 1.23 Range Iterator to stream messages
		for resp, err := range client.Models.GenerateContentStream(ctx, u.config.GeminiChatModel, contents, config) {
			if err != nil {
				_ = pw.CloseWithError(err)
				return
			}

			if resp != nil {
				// Use helper method .Text() to get the text
				_, _ = fmt.Fprint(pw, resp.Text())
			}
		}
	}()

	return pr, nil
}
