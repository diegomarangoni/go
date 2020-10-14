package foo

import "diegomarangoni.dev/go/pkg/lib/discovery"

var client *discovery.ClientDiscovery

func RegisterClient(cli *discovery.ClientDiscovery) {
	if client != nil {
		return
	}

	client = cli
}
