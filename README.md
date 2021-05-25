# github.release.notes

[![Go Report Card](https://goreportcard.com/badge/github.com/infobloxopen/atlas-cli)](https://goreportcard.com/report/github.com/infobloxopen/github.release.notes)

The initial version of this tool was created during **"Hack Week 2021"** Infoblox event

## Content

- [Content](#Content)
- [What is this?](#What-is-this)
- [Why github.release.notes?](#Why-github.release.notes)
- [Getting Started](#Getting-Started)
  - [Usage](#Usage)
    - [Github Actions Usage](#Github-Actions-Usage)
    - [Local Usage](#Local-Usage)
  - [Flags](#Flags)
  - [Customize your release notes](#Customize-your-release-notes)
  - [Versioning](#Versioning)

  ------------------------------------------------------------------------------------------

## What is this?

This tool allows to automate process of release notes generation and publishing. It collects all necessary information based on tags and commits from the repository, generates release notes data with Markdown formatting and publishes release notes.

### *Why do we need release notes?*

We need it to make it easier for users and contributors to see precisely what changes have been made between each release of the project.

### *Why should we care?*

Because software tools are for _people_. And people need to know about changes. Whether consumers or developers, the end users of software are human beings who care about what's in the software. When the software changes, people want to know why and how.

  ------------------------------------------------------------------------------------------

## Why github.release.notes?

There are many similar tools in the internet in general and at GitHub in particular. All of them may provide you a possibility to create release notes (some do this better, some not enough).
The main advantage of the **github.release.notes** tool are using of [**templates**](#Customize-your-release-notes) and well documented [**GitHub actions approach**](#Github-Actions-Usage)

  ------------------------------------------------------------------------------------------

## Getting Started

These instructions will get you a copy of Release-notes command-line tool up and running on your local machine or as a github actions

  ------------------------------------------------------------------------------------------

### Usage

There are several ways to use **github.release.notes**:

- [Release Notes Via Github Actions](#Github-Actions-Usage)
- [Release Notes Via Mannual Running](#Local-Usage)

#### Github Actions Usage

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
          template-title: templates/defaultTitle.tmpl
          template-body: templates/defaultBody.tmpl
```

[Back to content](#Content)

#### Local Usage

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

Also you may just pull docker image:

```sh
docker pull dlahuta/github.release.notes:v1.0.0
```

And then use it for creation of the release notes:
```
docker run --rm dlahuta/github.release.notes:v1.0.0 --github.token=<GITHUB_PAT> \
                                                    --github.repository=<GITHUB_REPOSITORY_NAME> \
                                                    --github.owner=<GITHUB_REPOSITORY_OWNER> \
                                                    --github.actor=<GITHUB_USER>
```

[Back to content](#Content)

  ------------------------------------------------------------------------------------------

### Flags

Here's the full set of flags for application.

| Flag                      | Description                                                             | Required      | Default Value         |
| ------------------------- | ----------------------------------------------------------------------- | ------------- | --------------------- |
| `app.id`                  | identifier for the application                                          | No            | `""`                  |
| `github.actor`            | github user name                                                        | No            | `""`                  |
| `github.owner`            | github owner (the same as `github.repository.owner`, use one of them)   | Yes           | `""`                  |
| `github.repository`       | github repository name                                                  | Yes           | `""`                  |
| `github.repository.owner` | github owner (the same as `github.owner`, use one of them)              | Yes           | `""`                  |
| `github.tag`              | github repository tag (all tags will be processed if skipped)           | No            | `""`                  |
| `github.token`            | github Personal Access Token for communication with gihub.api           | Yes           | `""`                  |
| `logging.level`           | log level of application                                                | No            | `"info"`              |
| `template.body`           | template file for release notes' body                                   | No            | `"defaultBody.tmpl"`  |
| `template.title`          | template file for release notes' title                                  | No            | `"defaultTitle.tmpl"` |
| `update.exist`            | is it required to update the release notes in case it already exists    | No            | `"false"`             |

Possible log levels: `"panic"`, `"fatal"`, `"error"`, `"warn"`, `"warning"`, `"info"`, `"debug"`, `"trace"`
**NOTE:** the `github.owner` and `github.repository.owner` both specify the same. Use one of them. Otherwise the last specified will be used.

[Back to content](#Content)

  ------------------------------------------------------------------------------------------

### Customize your release notes

One of the main advantages of the **github.release.notes** tool is using of **templates** for release notes.
So, no more hardcoded forms! Now you can decide what you want to see in your release notes and how.
Since the tool allows to specify template as an argument you may have several templates for different releases.

The tool uses [**golang template**](https://golang.org/pkg/text/template/) syntax to define templates for release notes.

Here is a simple example of the default release notes template:

- title:
  ```
  [{{.Tag}}] {{.TagComment}} ({{.TagDate}})
  ```
- body:
  ```
  {{if .ChangeLogLink -}}
  [Full Changelog]({{.ChangeLogLink}})
  
  {{end -}}
  **New commits and merged pull requests:**
  {{range .Commits}}
  - {{.CommitMessage}}
  {{- if .CommitPR}} {{.CommitPR}} {{end}}
  {{- if and (not .CommitPR) .CommitURL}} [[link]({{.CommitURL}})] {{end -}}
  ([{{.CommitAuthor}}]({{.CommitAuthorURL}})) {{.CommitDate}}
  {{- end}}
  ```

In order to use your template:
1. Create a file with your template (e.g. `/home/templates/myTemplateTitle.tmpl` and `/home/templates/myTemplateBody.tmpl`)
2. Use it while local execution via binary:
```
release-notes --github.token=<GITHUB_PAT> \
              --github.repository=<GITHUB_REPOSITORY_NAME> \
              --github.owner=<GITHUB_REPOSITORY_OWNER> \
              --template.title=/home/templates/myTemplateTitle.tmpl\
              --template.body=/home/templates/myTemplateBody.tmpl
```
3. To use it via docker approach you need to mount to docker container this file, e.g.
```
docker run --rm -v /home/templates:/templates dlahuta/github.release.notes:v1.0.0 \
                                                    --github.token=<GITHUB_PAT> \
                                                    --github.repository=<GITHUB_REPOSITORY_NAME> \
                                                    --github.owner=<GITHUB_REPOSITORY_OWNER> \
                                                    --template.title=/templates/myTemplateTitle.tmpl\
                                                    --template.body=/templates/myTemplateBody.tmpl
```

Here's the full set of values you may use in templates:
```
- Tag                 - name of the tag in the repo
- TagComment          - comment message added to the tag
- TagDate             - date of the tag creation
- ChangeLogLink       - link to the changelog between the current tag and previous one
- Commits             - list of the commits in the release
  - CommitAdditions   - total additions within the commit
  - CommitAuthor      - author of the commit
  - CommitAuthorURL   - HTML URL of the commit's author
  - CommitDate        - date of the commit creation
  - CommitDeletions   - total deletions within the commit
  - CommitMessage     - comment message of the commit
  - CommitPR          - pull request related to the commit (if exists)
  - CommitURL         - HTML URL of the commit (may be used if no PR related to the commit)
```

[Back to content](#Content)

  ------------------------------------------------------------------------------------------

### Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/infobloxopen/github.release.notes/tags).

[Back to content](#Content)

  ------------------------------------------------------------------------------------------
