# Blue Green Deployment of this App

## Pre-Reqs

- K8S Cluster with RBAC enabled
- [Ingress Controller](https://github.com/heptio/contour)
- [Argo Workflow Controller](https://argoproj.github.io/argo)
- GoLang 1.10+
- [Kustomize](https://github.com/kubernetes-sigs/kustomize)

### Install Kustomize

```bash
brew install kustomize
```

### Install Ingress Controller and Ingress Resource

```bash
kubectl apply -f https://j.hept.io/contour-deployment-rbac
```

If you are running k8s cluster with docker for mac make the service type as NodePort

```bash
kubectl -n heptio-contour patch svc contour --type json -p='[{"op":"replace","path":"/spec/type", "value": "NodePort"}]'
```

### Install Argo WF Controller

```bash
kubectl create ns argo
kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo/v2.2.1/manifests/install.yaml
```

### Install the GoReactAPP

#### Bootstrap

- Create a Namespace for Staging

```bash
kubectl create ns staging
```

- Create a Ingress Resource

```bash
kubectl -n staging apply -f deployments/ingress.yaml
```

- Deploy the first version

```
kustomize build deployments/overlays/staging/