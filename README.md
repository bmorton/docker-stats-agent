# Docker stats collector to Statsd

This agent will run against the configured Docker daemon and watch for stats on all running containers and push those stats to a Statsd collector.


## WARNING

This is a work in progress and doesn't actually do everything it claims to do yet.


## Configuration

The agent can be configured via CLI flags as well as environment variables that match (e.g. `-docker-cert-path` maps to `DOCKER_CERT_PATH`).

```ShellOutput
$ docker-statsd-agent --help
Usage of docker-statsd-agent:
  -docker-cert-path="": path to the cert.pem, key.pem, and ca.pem for authenticating to Docker
  -docker-host="unix:///var/run/docker.sock": address of Docker host
  -docker-tls-verify=false: use TLS client for Docker
```
