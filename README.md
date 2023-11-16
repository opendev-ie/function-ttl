# function-ttl
[![CI](https://github.com/crossplane-contrib/function-ttl/actions/workflows/ci.yml/badge.svg)](https://github.com/crossplane-contrib/function-ttl/actions/workflows/ci.yml)

A [composition function][functions] in [Go][go].

`function-ttl` allows you to set a time-to-live (ttl) on Composite Resources. This is done via an annotation: `fn.crossplane.io/ttl`. Once the specified duration has passed, on the next reconciliation loop, and every loop thereafter, the function will clear the desired set of composed resources, causing any already created to be removed.

This function should come towards the end of the pipeline, after any functions that may add composed resources to the desired state.

## Building Locally

This function is built using [Go][go], [Docker][docker], and the [Crossplane CLI][cli], and was generated from the [function template][function template]

```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests - see fn_test.go
$ go test ./...

# Build the function's runtime image - see Dockerfile
$ docker build . --quiet --platform=linux/amd64 --tag runtime-amd64
$ docker build . --quiet --platform=linux/arm64 --tag runtime-arm64

# Build a function package - see package/crossplane.yaml
$ crossplane xpkg build --package-root=package --embed-runtime-image=runtime-amd64 --package-file=function-amd64.xpkg
$ crossplane xpkg build --package-root=package --embed-runtime-image=runtime-arm64 --package-file=function-arm64.xpkg

# Push the function image
$ crossplane xpkg push --package-files=function-amd64.xpkg,function-arm64.xpkg {{registry}}/function-ttl:{{version}}
```

[functions]: https://docs.crossplane.io/latest/concepts/composition-functions
[go]: https://go.dev
[function guide]: https://docs.crossplane.io/knowledge-base/guides/write-a-composition-function-in-go
[function template]: https://github.com/crossplane/function-template-go
[package docs]: https://pkg.go.dev/github.com/crossplane/function-sdk-go
[docker]: https://www.docker.com
[cli]: https://docs.crossplane.io/latest/cli