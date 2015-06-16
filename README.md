# Docker Stats Agent

This agent will run against the configured Docker daemon and watch for stats on all running containers and push those stats to a metrics collector such as Statsd.


## WARNING

This is a work in progress and doesn't actually do everything it claims to do yet.


## Configuration

The agent can be configured via CLI flags as well as environment variables that match (e.g. `-docker-cert-path` maps to `DOCKER_CERT_PATH`).

```ShellOutput
$ docker-stats-agent --help
Usage of docker-stats-agent:
  -docker-cert-path="": path to the cert.pem, key.pem, and ca.pem for authenticating to Docker
  -docker-host="unix:///var/run/docker.sock": address of Docker host
  -docker-tls-verify=false: use TLS client for Docker
```


## Example

```ShellOutput
$ docker-stats-agent
Querying for running containers...
Waiting for stats...
Watching ba686b27f38d...
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 8.976 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 8.976 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 8.976 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 8.976 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 8.976 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 8.976 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 8.976 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 8.976 kB/648 B
Watching feca1ee0100e...
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 9.144 kB/648 B
feca1ee0100e  0.00% 974.8 kB/2.105 GB 0.05% 168 B/168 B
feca1ee0100e  0.00% 974.8 kB/2.105 GB 0.05% 418 B/418 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 9.394 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 9.484 kB/648 B
feca1ee0100e  0.00% 974.8 kB/2.105 GB 0.05% 508 B/508 B
feca1ee0100e  0.00% 950.3 kB/2.105 GB 0.05% 508 B/508 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 9.484 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 9.484 kB/648 B
feca1ee0100e  0.00% 950.3 kB/2.105 GB 0.05% 508 B/508 B
feca1ee0100e  0.37% 450.6 kB/2.105 GB 0.02% 578 B/578 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 9.554 kB/648 B
Closing feca1ee0100e...
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 9.554 kB/648 B
ba686b27f38d  0.00% 3.838 MB/2.105 GB 0.18% 9.554 kB/648 B
```
