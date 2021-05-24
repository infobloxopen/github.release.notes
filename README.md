# github.release.notes

[![Go Report Card](https://goreportcard.com/badge/github.com/infobloxopen/atlas-cli)](https://goreportcard.com/report/github.com/infobloxopen/github.release.notes)

_This generated README.md file loosely follows a [popular template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)._

One paragraph of project description goes here.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Installing

The following steps will install the `release-notes` binary to your `$GOBIN` directory.

```sh
$go get github.com/infobloxopen/github.release.notes/release-notes

```

You're all set! Alternatively, you can clone the repository and install the binary manually.

```sh
$git clone https://github.com/infobloxopen/github.release.notes.git
$cd github.release.notes
$make
```

#### Flags

Here's the full set of flags for application.

| Flag          | Description                                                         | Required      | Default Value |
| ------------- | ------------------------------------------------------------------- | ------------- | ------------- |
| `github.token`| The GitHub Personal Access Token for communication with gihub.api   | Yes           | `""`          |
| `github.tag`  | Github repository tag                                               | No            | `""`          |

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/infobloxopen/atlas-cli/github.release.notes/tags).
