---
title: "Kustomized Installation"
weight: 210
---

We are using Flux to package our Kustomize files through OCI, and they are built and released just as our helm solution.

There is no way of downloading manifest files through the [Kustomize CLI](https://kustomize.io/), so if you want to
download the Kustomize manifest you need to install the [Flux cli](https://fluxcd.io/flux/installation/).

## Install Flux Cli

Mac

```shell
brew install fluxcd/tap/flux
```

Linux

```shell
curl -s https://fluxcd.io/install.sh | sudo bash
```

## Download Kustomize files

After you have downloaded Flux you can use `flux pull artifact` to download the manifests.

This will provide you the manifest files unpacked and ready to use.

```shell
mkdir doris-operator
flux pull artifact oci://ghcr.io/linsoss/kustomize/doris-operator:1.0.0 --output ./doris-operator/
```

## Install Operator

```shell
kubectl apply -k doris-operator
```

## Uninstall Operator

```shell
kubectl delete -k doris-operator
```

