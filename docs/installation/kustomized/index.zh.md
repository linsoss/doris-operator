---
title: "Kustomized 方式"
weight: 210
---

Doris Operator使用 Flux 通过OCI打包 Kustomize文件，这和 Helm 的发布方式是一致的，因此如果您想要下载 Kustomize
文件清单，请先安装 [Flux cli](https://fluxcd.io/flux/installation/)。以下包含一个简单的 Flux Cli 安装步骤。

## 安装 Flux Cli

{{< tabs "install-flux" >}}
{{< tab "Mac" >}}

```shell
brew install fluxcd/tap/flux
```

{{< /tab >}}
{{< tab "Linux" >}}

```shell
curl -s https://fluxcd.io/install.sh | sudo bash
```

{{< /tab >}}
{{< /tabs >}}

## 下载 Kustomize文件

当你下载完 Flux 后，您可以使用 `flux pull artifact`  来下载 Kustomize 清单，这将为您提供已解压并准备好使用的清单文件，您可以根据您的需求修改其中的部分参数。

```shell
mkdir doris-operator
flux pull artifact oci://ghcr.io/linsoss/kustomize/doris-operator:1.0.1 --output ./doris-operator/
```

## 安装 Operator

```shell
kubectl apply -k doris-operator
```

## 卸载 Operator

```shell
kubectl delete -k doris-operator
```

