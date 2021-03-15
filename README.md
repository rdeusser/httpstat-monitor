# httpstat-monitor

A service that checks https://httpstat.us/200 and https://httpstat.us/503 to provide metrics on what the response time
is and if they're up.

## Testing

Run `make test`.

## Running locally

To build and run this project locally, you need the following things:

- Go v1.16
- Kubectl v1.18.x
- Helm v3.5.x
- Kubernetes cluster (e.g. Docker for Desktop with Kubernetes enabled)

To deploy everything:

1. `make deploy`
2. `make deploy-prometheus`

This will deploy httpstat-monitor and the [kube-prometheus stack](https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack) which includes the Prometheus
Operator, Prometheus, and Grafana (as well as a few other things).

If your Kubernetes cluster is running on your local machine, run `make port-forward-prometheus` to have kubectl
port-forward the Prometheus UI to your local machine and then navigate to `http://localhost:9090/graph` in your browser.

You can also do this with the Grafana UI by running `make port-forward-grafana` and then navigate to `http://localhost:8080`.

The default credentials for Grafana are:

User: `admin`
Password: `prom-operator`

## Screenshots

![Screenshot of Prometheus](/docs/screenshots/prometheus.png?raw=true)

![Screenshot of Grafana Dashboard](/docs/screenshots/grafana.png?raw=true)
