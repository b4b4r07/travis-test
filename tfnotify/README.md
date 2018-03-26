tfnotify
========

tfnotify parses the execution result of Terraform and notifies it to any Platform such as GitHub

## Motivation

There are commands such as plan and apply on Terraform command, but many developers think they would like to check if the execution of those commands succeeded.
Terraform commands are often executed via CI like Circle CI, but in that case you need to go to the CI page to check it.
This is very troublesome. It is very efficient if you can check it with GitHub's comment or Slack.
You can do this by using this command.

## Usage

This tool is still under alpha phase do the following usage may be changed without any notice.

```console
$ tfnotify --help
NAME:
   tfnotify - Notify the execution result of terraform command

USAGE:
   tfnotify [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     fmt      Parse stdin as a fmt result
     plan     Parse stdin as a plan result
     apply    Parse stdin as a apply result
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ci value        name of CI to run tfnotify
   --config value    config path (default: tfnotify.yaml)
   --notifier value  notification destination
   --help, -h        show help
   --version, -v     print the version
```

tfnotify accepts standard input. Also, it needs a configuration file. If `--config` is omitted, `tfnotify.yaml` in current working directory is loaded as its config file.

```console
$ terraform plan | tfnotify plan
```

Example of the config file is below.

```yaml
---
repository:
  owner: "kouzoh"
  name: "mercari-p-b4b4r07"
ci: circleci
notifier:
  github:
    token: $GITHUB_TOKEN
  slack:
    token: $SLACK_TOKEN
    channel: C7KCTQ734
terraform:
  plan:
    template: |
      {{ .Title }}
      {{ .Message }}
      {{if .Result}}
      <pre><code> {{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>
      <pre><code> {{ .Body }}
      </pre></code></details>
  apply:
    template: |
      {{ .Title }}
      {{ .Message }}
      {{if .Result}}
      <pre><code> {{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>
      <pre><code> {{ .Body }}
      </pre></code></details>
```

There is no need to replace TOKEN string such as `$GITHUB_TOKEN` with the actual token. Instead, it must be defined as an environment variable.

## Installation

```console
$ git clone https://github.com/kouzoh/tfnotify; cd tfnotify
$ dep ensure
$ go install
```

## License

MIT

## Auther

b4b4r07

