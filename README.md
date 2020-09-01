# promcollectors

Prometheus utility collectors that could not be found in official lib, like a counter in which you can issue a "set" value.

The only component here is the settablecounter now. Useful when getting counters from another system. Use wisely because counters are meant to be always going UP and go to zero when resetting it. Don't use this for gauge data!

Thanks https://www.robustperception.io/setting-a-prometheus-counter

## Usage

```golang

    import "github.com/flaviostutz/promcollectors"
...

	hostInfo = promcollectors.NewSettableCounterVec(prometheus.Opts{
		Name: "mongo_server_uptime_seconds",
		Help: "Basic server info and uptime in seconds",
	}, []string{
		"host",
		"version",
		"process",
    })
    
	prometheus.MustRegister(hostInfo)

    hostInfo.Set(6534, "host1", "4.4.0", "mongos")
...
```

See an example usage at http://github.com/stutzlab/mongo-simple-exporter/server_status.go

