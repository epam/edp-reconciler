package perfserver

import perfApi "github.com/epam/edp-perf-operator/v2/pkg/apis/edp/v1"

type PerfServer struct {
	Name      string
	Available bool
}

func ConvertPerfServerToDto(server perfApi.PerfServer) PerfServer {
	return PerfServer{
		Name:      server.Name,
		Available: server.Status.Available,
	}
}
