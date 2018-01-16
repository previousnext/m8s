Gopher
======

[![CircleCI](https://circleci.com/gh/previousnext/gopher.svg?style=svg)](https://circleci.com/gh/previousnext/gopher)

![Logo](/logo/small.png "Logo")

**Maintainer**: Gopher

This is a brief description on what the project does.

## Development

### Getting started

To work on this project you will first need Go installed on your machine.

#### Manual Setup

First make sure Go is properly installed and that a GOPATH has been set. You will also need to add $GOPATH/bin to your $PATH. For steps on getting started with Go: https://golang.org/doc/install

Next, using Git, clone this repository into $GOPATH/src/github.com/previousnext/gopher. All the necessary dependencies are either vendored or automatically installed, so you just need to type `make test`. This will run the tests and compile the binary. If this exits with exit status 0, then everything is working!

```bash
$ cd "$GOPATH/src/github.com/previousnext/gopher"
$ make test
```

To compile a development version of gopher, run `make build`. This will build everything using gox and put binaries in the bin and $GOPATH/bin folders:

```bash
$ make build
...

# Linux:
$ bin/gopher_linux_amd64 --help

# OSX:
$ bin/gopher_darwin_amd64 --help
```

#### Easy Setup

Alternatively, you can use the [Docker Compose](docker-compose.yml) stack in the root of this repo to stand up a container with the appropriate dev tooling already set up for you.

Using Git, clone this repo on your local machine. Run the test suite to ensure the tooling works.

```bash
$ docker-compose run --rm dev make test
```

To compile a development version of gopher, run `make build`. This will build everything using gox and put binaries in the bin and $GOPATH/bin folders:

```bash
$ docker-compose run --rm dev make build

...

$ docker-compose run --rm dev bin/gopher_linux_amd64 --help
```

### Dependencies

gopher stores its dependencies under `vendor/`, which [Go 1.6+ will automatically recognize and load](https://golang.org/cmd/go/#hdr-Vendor_Directories). We use [`dep`](https://github.com/golang/dep) to manage the vendored dependencies.

If you're developing m8s, there are a few tasks you might need to perform.

For details, see:

* [Adding a dependency](#adding-a-dependency)
* [Updating a dependency](#updating-a-dependency)

### Documentation

See `/docs`

### Resources

* [Dave Cheney - Reproducible Builds](https://www.youtube.com/watch?v=c3dW80eO88I)
* [Bryan Cantril - Debugging under fire](https://www.youtube.com/watch?v=30jNsCVLpAE&t=2675s)
* [Sam Boyer - The New Era of Go Package Management](https://www.youtube.com/watch?v=5LtMb090AZI)
* [Kelsey Hightower - From development to production](https://www.youtube.com/watch?v=XL9CQobFB8I&t=787s)

### Tooling

```bash
# Dependency management
go get -u github.com/golang/dep/cmd/dep

# Testing
go get -u github.com/golang/lint/golint

# Release management.
go get -u github.com/tcnksm/ghr

# Build
go get -u github.com/mitchellh/gox
```

### Common Tasks

#### Adding a dependency

If you're adding a dependency, you'll need to vendor it in the same Pull Request as the code that depends on it. You should do this in a separate commit from your code, as makes PR review easier and Git history simpler to read in the future.

To add a dependency:

Assuming your work is on a branch called `my-feature-branch`, the steps look like this:

1. Vendor the new dependency.

    ```bash
    dep ensure -add github.com/foo/bar
    ```

2. Review the changes in git and commit them.

#### Updating a dependency

To update a dependency:

1. Update the dependency.

    ```bash
    dep ensure -update github.com/foo/bar
    ```

2. Review the changes in git and commit them.

#### Running quality checks

```bash
make lint test
```

#### Building binaries

```bash
make build
```

#### Release

Release artifacts are pushed to the [github releases page](https://github.com/previousnext/gopher/releases) when tagged
properly. Use [semantic versioning](http://semver.org/) prefixed with `v` for version scheme. Examples:

- `v1.0.0`
- `v1.1.0-beta1`
