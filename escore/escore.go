package escore

import "github.com/elastic/go-elasticsearch/v9"

type ESConfig struct {
	Address  []string
	Username string
	Password string
}

func Connect(conf *ESConfig) (*elasticsearch.TypedClient, error) {
	esClient, err := elasticsearch.NewTyped(
		elasticsearch.WithAddresses(conf.Address...),
		elasticsearch.WithBasicAuth(conf.Username, conf.Password),
		elasticsearch.WithLogger(&Logger{}),
	)
	if err != nil {
		return nil, err
	}
	return esClient, nil
}
