package config

type ConfigServer struct {
	secret string
}

func (cs *ConfigServer) Secret() string {
	return cs.secret
}
