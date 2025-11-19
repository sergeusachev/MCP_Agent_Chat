package gigachat

import (
	"os"
	"strings"
)

func GetOauthToken() string {
	return ReadTextFromFile("/Users/sergeyusachev/Projects/GoProjects/MCP_Example/api/gigachat/secret/oauth_gigachat_token.txt")
}

func ReadTextFromFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(data))
}
