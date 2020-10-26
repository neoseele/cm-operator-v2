# cm-operator-v2

Rebuild [cm-operator](https://github.com/neoseele/cm-operator) with [operator-sdk](https://github.com/operator-framework/operator-sdk)

## Build

```sh
make clean
make build-dockerhub
```

## Run locally for debug

```sh
operator-sdk run --local --namespace=default
```

## Deploy

```sh
make cr
```

## Teardown

```sh
make teardown-crd
```

## Usage

### Annotate the pod that needs to be scraped

* `cm.example.com/scrape` (required)
* `cm.example.com/port` (optional, default: `80`)
* `cm.example.com/path` (optional, default: `/metrics`)

Example:

```sh
POD=some_pod
# create
kubectl annotate --overwrite pods $POD 'cm.example.com/scrape'='true' 'cm.example.com/port'='9990'
# remove
kubectl annotate --overwrite pods $POD 'cm.example.com/scrape-' 'cm.example.com/port-'
```

### Annotate the node that needs to be scrapes

> port/path is hardcoded to cadvisor endpoint `:10255/metrics/cadvisor`

* `cm.example.com/scrape` (required)

Example:

```sh
NODE=some_node
# create
kubectl annotate --overwrite nodes $NODE 'cm.example.com/scrape'='true'
# remove
kubectl annotate --overwrite nodes $NODE 'cm.example.com/scrape-'
```

### Create the CR

The listed metrics will be sent to Cloud Monitoring

```yaml
apiVersion: cm.example.com/v1alpha1
kind: CustomMetric
metadata:
  name: cm
spec:
  project: nmiu-play
  cluster: ebpf
  location: australia-southeast1-a
  metrics:
    - cilium_*
    - container_network_*
```
