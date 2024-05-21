package lib

const (
	DefaultEpisodesAmount = 10
	DefaultBaseUrl        = "https://rickandmortyapi.com/api"
	DefaultTimeout        = 5
	DefaultCharacterPath  = "/character"
	DefaultEpisodePath    = "/episode"
)

type ClientSettings struct {
	Timeout int    `json:"timeout"`
	BaseUrl string `json:"baseUrl"`
}

type ServiceSettings struct {
	EpisodesAmount int    `json:"episodesAmount"`
	CharacterPath  string `json:"characterPath"`
	EpisodePath    string `json:"episodePath"`
}

type Config struct {
	ServiceConfig ServiceSettings `json:"serviceConfig"`
	ClientConfig  ClientSettings  `json:"clientConfig"`
}

func LoadConfig() Config {
	return Config{
		ServiceConfig: ServiceSettings{
			EpisodesAmount: DefaultEpisodesAmount,
			CharacterPath:  DefaultCharacterPath,
			EpisodePath:    DefaultEpisodePath,
		},
		ClientConfig: ClientSettings{
			Timeout: DefaultTimeout,
			BaseUrl: DefaultBaseUrl,
		},
	}
}
