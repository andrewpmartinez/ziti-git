# ziti-git
A tool to manage multiple repositories with special considerations for the github.com/openziti project

## Requirements
- git >= 1.14

## Installation
```
go get -u github.com/andrewpmartinez/ziti-git
go install github.com/andrewpmartinez/ziti-git
```

## Alias to `zg`
```
echo 'alias zg=$GOPATH/bin/ziti-git' >> ~/.bashrc
```

## Usage
```
Ziti Git is a multi-repo git tool with additions for the open ziti project!

Usage:
  ziti-git [flags]
  ziti-git [command]

Available Commands:
  branch         list all repo branches or repos in <tag>
  clone          clones the core openziti repos to the current directory
  help           Help about any command
  list           list all repos or repos for <tag>
  register       add the repo in <path> to the list of repos, with an optional <tag>
  table-status   show the table status of all the repos or of a specific tag
  unregister     unregister <repo>
  unregister-tag unregister-tag <tag>

Flags:
  -h, --help         help for ziti-git
  -t, --tag string   limits actions to repos with <tag>

Use "ziti-git [command] --help" for more information about a command.
```

## Getting Started

To start hacking away on Ziti first clone the repos:

```
mkdir ziti
cd ziti
ziti-git clone
```

If you want to add these new repos to ziti-gits repo list use the `-r` flag
and optionally the `-t` flag to set a specific tag.

```
mkdir ziti
cd ziti
ziti-git clone -r -t myTag
```

## Table Status

A tabular Git status can be displayed by using the `table-status` or `ts` command. The output
can be limited by specifying a specific tag via `-t`.

```
> ziti-git table-status
+------------+--------------+----------+--------+----------+-------------------------------------------+
|    NAME    |    BRANCH    |   TAG    | STAGED | UNSTAGED |                 LOCATION                  |
+------------+--------------+----------+--------+----------+-------------------------------------------+
| foundation | master       | v0.12.0  |        |          | /home/user/repos/openziti/foundation      |
| ziti       | release-next | 3a19537  |        |          | /home/user/repos/openziti/ziti            |
| edge       | master       | v0.15.40 |        |          | /home/user/repos/openziti/edge            |
| fabric     | master       | v0.12.1  |        |          | /home/user/repos/openziti/fabric          |
+------------+--------------+----------+--------+----------+-------------------------------------------+
```

## Fetching On All Repos

Arbitrary git command can be executed on the entire set of repos or
sets defined by tags. In this example `git fetch` will be executed on
all repos.

```
> ziti-git fetch
```

Or on a specific tag:

```
> ziti-git -t myTag fetch
```

This can also be used to create branches, checkout branches, hard
reset, etc. across all repos.

## Unregister Repo

Repositories can be removed by location or by tag. To remove a specific
repository by path:

```
> ziti-git unregister ./edge
```

To remove all repositories with a specific tag:
```
> ziti-git unregister-tag myTag
```

## Aliases

Most ziti-git commands have short aliases:

```
  b  = branch 
  c  = clone
  l  = list
  r  = register
  ts = table-status
  u  = unregister
  ut = unregister-tag

```

Aliases can be found by use the `-h` flag on commands in the "Aliases" section:

```
> ziti-git register -h
add the repo in <path> to the list of repos, with an optional <tag>

Usage:
  ziti-git register [-t <tag>] <path> [flags]

Aliases:
  register, r

Flags:
  -h, --help   help for register

Global Flags:
  -t, --tag string   limits actions to repos with <tag>
```


## Acknowledgements
Ziti Git is based off of [gmg](https://github.com/abrochard/go-many-git) which in turn was inspired by the amazing [mr](https://myrepos.branchable.com) and [gr](https://github.com/mixu/gr) tools.

A big thanks to all.
