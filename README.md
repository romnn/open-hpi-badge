## open-hpi-badge

[![Build Status](https://travis-ci.com/romnnn/open-hpi-badge.svg?branch=master)](https://travis-ci.com/romnnn/open-hpi-badge)
[![GitHub](https://img.shields.io/github/license/romnnn/open-hpi-badge)](https://github.com/romnnn/open-hpi-badge)
[![GoDoc](https://godoc.org/github.com/romnnn/open-hpi-badge?status.svg)](https://godoc.org/github.com/romnnn/open-hpi-badge) [![Docker Pulls](https://img.shields.io/docker/pulls/romnn/open-hpi-badge)](https://hub.docker.com/r/romnn/open-hpi-badge) [![Test Coverage](https://codecov.io/gh/romnnn/open-hpi-badge/branch/master/graph/badge.svg)](https://codecov.io/gh/romnnn/open-hpi-badge)
[![Release](https://img.shields.io/github/release/romnnn/open-hpi-badge)](https://github.com/romnnn/open-hpi-badge/releases/latest)

Tiny go server that serves a custom api endpoint for [shields.io](https://img.shields.io) to embed openHPI related badges on the web.

```bash
go get github.com/romnnn/open-hpi-badge/cmd/openhpibadge
```

You can also download pre built binaries from the [releases page](https://github.com/romnnn/open-hpi-badge/releases), or use the `docker` image:

```bash
docker pull romnn/openhpibadge
```

For a list of options, run with `--help`.

#### Deploying the server
There are a few options for deploying your endpoint:

1. Running the binary 
    ```bash
    go build github.com/romnnn/open-hpi-badge/cmd/openhpibadge
    ./openhpibadge --port 8080 --prod
    ```

    You can also download pre built binaries from the [releases page](https://github.com/romnnn/open-hpi-badge/releases)

2. Using `docker`
    ```bash
    docker run -p 8080:8080 romnn/open-hpi-badge --port 8080 --prod
    ```

#### Using the badges
https://img.shields.io/endpoint?url=...&style=...

#### Usage as a library

The package can also be imported as a library that exports the core functionality to build your more custom endpoint or service.

```golang
import github.com/romnnn/open-hpi-badge
```

For example, you can scrape an openHPI mooc by URL:

```golang
course, err := openhpibadge.ScrapeMOOCByURL("https://open.hpi.de/courses/neuralnets2020")
if err != nil {
    panic(err)
}
fmt.Println(course.Participants.Current)
```

```golang
course, err := openhpibadge.ScrapeMOOCByName("neuralnets2020")
if err != nil {
    panic(err)
}
fmt.Println(course.Participants.Current)
```

For more examples, see `examples/`.


#### Development

######  Prerequisites

Before you get started, make sure you have installed the following tools::

    $ python3 -m pip install -U cookiecutter>=1.4.0
    $ python3 -m pip install pre-commit bump2version invoke ruamel.yaml halo
    $ go get -u golang.org/x/tools/cmd/goimports
    $ go get -u golang.org/x/lint/golint
    $ go get -u github.com/fzipp/gocyclo
    $ go get -u github.com/mitchellh/gox  # if you want to test building on different architectures

**Remember**: To be able to excecute the tools downloaded with `go get`, 
make sure to include `$GOPATH/bin` in your `$PATH`.
If `echo $GOPATH` does not give you a path make sure to run
(`export GOPATH="$HOME/go"` to set it). In order for your changes to persist, 
do not forget to add these to your shells `.bashrc`.

With the tools in place, it is strongly advised to install the git commit hooks to make sure checks are passing in CI:
```bash
invoke install-hooks
```

You can check if all checks pass at any time:
```bash
invoke pre-commit
```

Note for Maintainers: After merging changes, tag your commits with a new version and push to GitHub to create a release:
```bash
bump2version (major | minor | patch)
git push --follow-tags
```

###### Internationalization

Developers who want to add or update translations can follow these steps:
```bash
cd cmd/openhpibadge
# Update the english ones as you wish and eventually sync them
goi18n extract -outdir intn
# This creates a translate.*.toml with all untranslated messages
goi18n merge -outdir intn intn/active.*.toml
# Merges the translate.*.toml back into the active.*.toml
goi18n merge -outdir intn intn/active.*.toml intn/translate.*.toml
# Stores the files as binary assets
go-bindata intn/
```

For information on internationalization, see [the guide](https://github.com/nicksnyder/go-i18n).

#### Note

This project is still in the alpha stage and should not be considered production ready.
