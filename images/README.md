# Docker files for Doris components

Base on https://github.com/apache/doris/tree/master/docker and optimize the entrypoint of FE, BE and make the resulting
image support multiple architectures.

### Build Steps

1. Make sure docker is installed and buildx is enabled.

2. Download the (Doris binary package)[https://doris.apache.org/download/] to the `dockerfiles/resource` directory, including x64 and arm versions, and unzip it with the following names:

   ```bash
   resource
   ├── apache-doris-bin-arm64
   └── apache-doris-bin-x86
   ```

3. Create a cross-platform docker builder for amd64 and arm64:

   ```bash
   docker buildx create --name doris-builder --use --platform linux/amd64,linux/arm64
   ```

4. Build cross-platform arch docker image:

   ```bash
   cd dockerfiles
   docker buildx build --platform linux/amd64,linux/arm64 -t <register>/<namespace>/doris-fe:2.0 -f fe/Dockerfile . --push
   docker buildx build --platform linux/amd64,linux/arm64 -t <register>/<namespace>/doris-be:2.0 -f be/Dockerfile . --push
   docker buildx build --platform linux/amd64,linux/arm64 -t <register>/<namespace>/doris-cn:2.0 -f cn/Dockerfile . --push
   docker buildx build --platform linux/amd64,linux/arm64 -t <register>/<namespace>/doris-broker:2.0 -f broker/Dockerfile . --push
   ```
