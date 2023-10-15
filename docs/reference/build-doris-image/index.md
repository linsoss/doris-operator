---
title: "Build the Doris image"
weight: 710
---

The Doris Operator uses a special cross-architecture Doris image and is not compatible with
the [apache/doris](https://hub.docker.com/r/apache/doris) image on Docker Hub.

## Environment

- Docker version 19.03+
- Enable Docker [buildx](https://github.com/docker/buildx) plugin

## Steps

1. Clone the Doris repository and enter the `images` directory:

   ```bash
   git clone https://github.com/linsoss/doris-operator
   cd doris-operator/images
   ```

2. Download the [Doris binary package](https://doris.apache.org/download/) to the `dockerfiles/resource` directory,
   including x64 and arm versions, and unzip it with the following names:

   ```other
   images
    ├── ...
    ├── resource
         ├── apache-doris-bin-arm64
         |     └── ...
         └── apache-doris-bin-amd64
               └── ...
   ```
3. Create a cross-platform docker builder for amd64 and arm64:

   ```bash
   docker buildx create --name doris-builder --use --platform linux/amd64,linux/arm64
   ```

4. Build cross-architecture Docker images and push them to the registry

   ```bash
   docker buildx build --platform linux/amd64,linux/arm64 -f fe/Dockerfile -t <register>/<namespace>/doris-fe:<tag> . --push
   docker buildx build --platform linux/amd64,linux/arm64 -f be/Dockerfile -t <register>/<namespace>/doris-be:<tag> . --push
   docker buildx build --platform linux/amd64,linux/arm64 -f cn/Dockerfile -t <register>/<namespace>/doris-cn:<tag> . --push
   docker buildx build --platform linux/amd64,linux/arm64 -f broker/Dockerfile -t <register>/<namespace>/doris-broker:<tag> . --push
   ```
