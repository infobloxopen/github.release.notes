# github.release.notes

[![Go Report Card](https://goreportcard.com/badge/github.com/infobloxopen/atlas-cli)](https://goreportcard.com/report/github.com/infobloxopen/github.release.notes)

## Getting Started

These instructions will get you a copy of Release-notes command-line tool up and running on your local machine or as a github actions

## Usage

You can use github.release.notes package to publish release notes.

- [Release Notes Via Github Actions](#Github-Actions-Usage)
- [Release Notes Via Mannual Running](#Local-Usage)

## Github Actions Usage

See [action.yml](action.yml)

Basic:

```yaml
name: release-notes
on:
  push:
    tags: ['v*']
  workflow_dispatch:
jobs:
  release-notes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set env
        id: vars        
        run: |
          echo ::set-output name=tag::${GITHUB_REF#refs/*/}
      - name: Run release.notes image
        uses: infobloxopen/github.release.notes@v1.0.0
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
          github-tag: ${{ steps.vars.outputs.tag }}
```

With Release Notes Template:

```yaml
name: release-notes
on:
  push:
    tags: ['v*']
  workflow_dispatch:
jobs:
  release-notes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set env
        id: vars        
        run: |
          echo ::set-output name=tag::${GITHUB_REF#refs/*/}
      - name: Run release.notes image
        uses: infobloxopen/github.release.notes@v1.0.0
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
          github-tag: ${{ steps.vars.outputs.tag }}
```

### Local Usage

The following steps will install the `release-notes` binary to your `$GOBIN` directory.

```sh
go get github.com/infobloxopen/github.release.notes/release-notes

```

You're all set! Alternatively, you can clone the repository and install the binary manually.

```sh
git clone https://github.com/infobloxopen/github.release.notes.git
cd github.release.notes
make
```

After installing please use the following command:

```sh
release-notes --github.token=<GITHUB_PAT> \
              --github.repository=<GITHUB_REPOSITORY_NAME> \
              --github.owner=<GITHUB_REPOSITORY_OWNER> \
              --github.actor=<GITHUB_USER>
```

#### Flags

Here's the full set of flags for application.

| Flag          | Description                                                         | Required      | Default Value |
| ------------- | ------------------------------------------------------------------- | ------------- | ------------- |
| `github.token`| The GitHub Personal Access Token for communication with gihub.api   | Yes           | `""`          |
| `github.tag`  | Github repository tag                                               | No            | `""`          |

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/infobloxopen/github.release.notes/tags).
