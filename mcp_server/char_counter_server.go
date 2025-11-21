package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TextToUnicodeInput struct {
	Text string `json:"text" jsonschema:"Text to convert to unicode codes"`
}

type TextToUnicodeOutput struct {
	UnicodeCodes string `json:"unicode_codes" jsonschema:"Unicode codes of the text characters"`
}

func TextToUnicode(ctx context.Context, req *mcp.CallToolRequest, input TextToUnicodeInput) (
	*mcp.CallToolResult,
	TextToUnicodeOutput,
	error,
) {
	var unicodeCodes []string
	for _, r := range input.Text {
		unicodeCodes = append(unicodeCodes, fmt.Sprintf("U+%04X", r))
	}
	result := strings.Join(unicodeCodes, " ")
	return nil, TextToUnicodeOutput{UnicodeCodes: result}, nil
}

func main() {
	// Create a server with text to unicode conversion tool
	server := mcp.NewServer(&mcp.Implementation{Name: "unicode-converter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "text_to_unicode", Description: "Converts text to unicode codes. Returns unicode code points for each character in the format U+XXXX."}, TextToUnicode)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
