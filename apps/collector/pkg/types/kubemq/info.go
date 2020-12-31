package kubemq

type Info struct {
	Host           string `json:"host"`
	Version        string `json:"version"`
	Mode           string `json:"mode"`
	IsHealthy      bool   `json:"is_healthy"`
	IsReady        bool   `json:"is_ready"`
	LeadershipRole string `json:"leadership_role"`
}

func NewInfo() *Info {
	return &Info{
		Host:           "",
		Version:        "",
		Mode:           "",
		IsHealthy:      false,
		IsReady:        false,
		LeadershipRole: "",
	}
}

func (i *Info) SetHost(value string) *Info {
	i.Host = value
	return i
}

func (i *Info) SetVersion(value string) *Info {
	i.Version = value
	return i
}

func (i *Info) SetMode(value string) *Info {
	i.Mode = value
	return i
}

func (i *Info) SetIsHealthy(value bool) *Info {
	i.IsHealthy = value
	return i
}

func (i *Info) SetIsReady(value bool) *Info {
	i.IsReady = value
	return i
}

func (i *Info) SetLeadershipRole(value string) *Info {
	i.LeadershipRole = value
	return i
}
