package cyoa

import (
	"encoding/json"
)

type Option struct {
	Chapter string `json:"arc"`
	Text    string `json:"text"`
}

type StoryArc struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Story map[string]StoryArc

func ParseStory(data []byte) (Story, error) {
	var story Story
	if err := json.Unmarshal(data, &story); err != nil {
		return nil, err
	}
	return story, nil
}
