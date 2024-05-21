package types

type Origin struct {
	Name string
	URL  string
}

type Location struct {
	Name string
	URL  string
}

type Character struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Species  string   `json:"species"`
	Type     string   `json:"type"`
	Gender   string   `json:"gender"`
	Origin   Origin   `json:"origin"`
	Location Location `json:"location"`
	Image    string   `json:"image"`
	Episode  []string `json:"episode"`
	URL      string   `json:"url"`
	Created  string   `json:"created"`
}

type Episode struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	Air_Date            string   `json:"air_date"`
	Episode             string   `json:"episode"`
	Characters          []string `json:"characters"`
	URL                 string   `json:"url"`
	Created             string   `json:"created"`
	SecondLastCharacter string
}
