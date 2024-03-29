# ziti-git

A tool to manage multiple repositories with special considerations for
the github.com/openziti project

## Requirements

- git >= 1.17

## Installation

### As library

`ziti-git` is developed as a library first and a tool second. If you wish to use `ziti-git`'s capabilities in your own
tooling simply add it as a dependency via `go get`

```
go get -u github.com/andrewpmartinez/ziti-git
```

### Build From Source Yourself

`ziti-git` has a builtin command/binary/executable called `ziti-git` that can be installed directly from your CLI as 
long as you  have golang installed. There are a few different flavors of installation you can choose. Below are examples
of installing `ziti-git` with go 1.17.

Install all commands:

`go install github.com/andrewpmartinez/ziti-git/...@latest`

Install command in the `cmd` dir:

`go install github.com/andrewpmartinez/ziti-git/cmd/...@latest`

Install only the `ziti-git` command:

`go install github.com/andrewpmartinez/ziti-git/cmd/ziti-git@latest`

As of now they are all equivalent as there is only one command in the project. Earlier versions of go may require
different syntax or commands.

### Checkout & Build

Alternatively if you want a local copy of the source code and build it you can checkout the repo and issue a build
command from within the repository folder:

```
git clone ...
cd ziti-git
go install ./...
```

## Alias to `zg`

For Linux

```
echo 'alias zg=$GOPATH/bin/ziti-git' >> ~/.bashrc
```

For Windows CMD:

```
doskey zg=ziti-git $*
```

For Powershell, add to the end of your $PROFILE:

```
Set-Alias zg ziti-git
```

## Aliases

Most ziti-git commands have short aliases. Setting up `ziti-git` as the
alias `zg` and using the command aliases can shorten the typing
necessary for repetitive tasks.

#### Example w/o aliases:

```
> ziti-git table-status
```

#### Example w/ aliases:

```
> zg ts
```

Here is a list of some aliases:

```
  e  = execute
  b  = branch 
  c  = clone
  g  = git
  l  = list
  r  = register
  ts = table-status
  u  = unregister
  ut = unregister-tag
  ul = use-local
  ur = use-remote
```

Aliases can be found by use the `-h` flag on commands in the "Aliases"
section:

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


## Usage

```
Ziti Git is a multi-repo git tool with additions for the open ziti project!

Usage:
  ziti-git [command]

Available Commands:
  branch         list all repo branches or repos in <tag>
  checkout       inspects the go.mod file of the openziti/ziti repo to produce a script to checkout exact openziti dependencies necessary
  clone          clones the core openziti repos to the current directory
  execute        execute commands across all repositories or specific <tag> repositories
  git            execute git commands across all repositories or specific <tag> repositories
  help           Help about any command
  list           list all repos or repos for <tag>
  register       add the repo in <path> to the list of repos, with an optional <tag>
  table-status   show the table status of all the repos or of a specific tag
  unregister     unregister <repo>
  unregister-tag unregister-tag <tag>
  use-local      alter go.mod files for ziti repos to use local repositories via replace directives
  use-remote     short cut for use-local -u

Flags:
  -h, --help         help for ziti-git
  -t, --tag string   limits actions to repos with <tag>

Use "ziti-git [command] --help" for more information about a command.
```

## Cloning -- Getting Started With Ziti

To start hacking away on Ziti first clone the `github.com/openziti/ziti`
repositories. It is suggested to run the `ziti-git clone` command inside
an empty directory as multiple directories will be created.

Example:

```
mkdir ziti
cd ziti
ziti-git clone
```

For easier management later, it is useful to register the cloned
repositories with `ziti-git` and specify a tag. This will make it easier
to manipulate them individually with the `-t` flag that is available on
most `ziti-git` commands.

```
mkdir ziti
cd ziti
ziti-git clone -r -t myZiti
```

You can clone then build to get your own copy of Ziti built and ready
for use:

```
mkdir myziti
cd myziti
ziti-git clone -r -t myZiti
cd ziti
go build ./...
```

The above will checkout the necessary Ziti repositories and then build
the Ziti binaries. They will end up `~/go/bin` or to the environment
variable path defined by `GOBIN` if set. The repository in the `ziti`
folder will contain the `openziti/ziti` repository which holds the code
that will build all of openziti's binaries.

## Table Status

A tabular Git status can be displayed by using the `table-status` or
`ts` command. The output can be limited by specifying a specific tag via
`-t`.

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

## Fetching On All Repositories / Arbitrary Git Commands

Arbitrary `git` command can be executed on the entire set of
repositories or sets defined by tags. In this example `git fetch` will
be executed on all repositories.

```
> ziti-git git fetch
```

Or the shorter command:

```
> zg g fetch
```

Or on a specific tag:

```
> ziti-git g -t myTag fetch
```

This can also be used to create branches, checkout branches, hard reset,
etc. across all repositories.

## Unregistering Repositories

Repositories can be removed by location or by tag. To remove a specific
repository by path:

```
> ziti-git unregister ./edge
```

To remove all repositories with a specific tag:

```
> ziti-git unregister-tag myTag
```

## Using Local -- Local Development

By default, building against the `openziti/ziti` repository folder
`ziti` will use its `go.mod` file to look up the correct versions to
build. If you would like to use only the locally checked out versions
(useful for developing locally) the `ziti-git use-local` command is
useful to update the `go.mod` file to add `replace` directives to use
your locally checked out versions. `go mod tidy` will be automatically
run unless the `--no-tidy` or `-n` is specified

The command makes the following assumptions:

- it is being run in the directory containing the `openziti/*`
  repositories
- it assumes that the `ziti`, `foundation`, `edge`, and `fabric` folders
  are siblings in said folder

```
> ziti-git use-local
```

Using the `use-local` command will alter the `go.mod` file across some
or all of the repositories mentioned above (depending on usage).
Committing modified `go.mod` files with `replace` directives is
generally not advised unless it is for your own personal use.

To reverse this process use:

```
> ziti-git use-local --undo
```

To limit the scope of `use-local` the `--current` flag can be used
within a specific repository folder to alter only the `go.mod` folder of
that repository.

```
> cd edge
> ziti-git use-local --current
```

`--current` can also be combined with `--undo` to limit the undo to only
the current repository.

```
> cd edge
> ziti-git use-local --current
> ziti-git use-local --current --undo
```

Specific repositories can also be swapped to use the locally checked out
versions by specifying them via the `--repo` flag.

The following would only use the local `edge` repository.

```
> ziti-git use-local --repo .*?edge.*?
```

Note that the repo flag treats the input as a regular expression. If
your shell requires escape characters (i.e. bash and `\`) those must be
applied on top of any regular expression escaping necessary.

The `--repo` flag can also be combined with `--current` and `--undo`. It
may also be specified multiple times.


## Using Remote -- Undoing Use Local

The command `ziti-git use-remote` can be used as a shortcut to using
`ziti-git use-local -u`. Both work exactly the same - the only
difference is that `use-remote` does not have a `-u` flag. `go mod tidy` 
will be automatically run unless the `--no-tidy` or `-n` is specified.

## Shortcut Repository Names

For Ziti specific repos, there are a number of shortcuts so that the full
regular expression does not have to be used for `use-local` and `use-remote`.

For example instead of this command:

```
> ziti-git use-local --repo .*?edge.*?
```

The `edge` shortcut can e used instead

```
> ziti-git use-local --repo edge
```

Here is a full list of the shortcuts defined:

- `edge`
- `fabric`
- `foundation`
- `sdk-golang`
- `channel`

### Checking Out Exact Matching Versions

When debugging issues or recreating historical versions, it is useful to
checkout the exact repository commits that were used to build a specific
version. The `ziti-git checkout` command can do that for you.

If you wish to checkout the commits used to build the `v0.16.0` of Ziti,
you can do the following:

```
> mkdir ziti-0.16.0
> cd ziti-0.16.0
> ziti-git clone -r -t v0.16.0
> cd ziti
> git checkout v0.16.0
> ziti-git checkout
```

Would output:

```
cd "/home/user/repos/ziti-0.16.0"
git -C "./edge" checkout v0.15.40
git -C "./fabric" checkout v0.12.1
git -C "./foundation" checkout v0.12.0
git -C "./sdk-golang" checkout v0.13.30
```

That output can be executed to checkout the proper versions. After that
the `use-local` command can be used to work on that specific version of
the openziti project - potentially to work on bug fix!

These repositories can then later be removed from `ziti-git` as the
`v0.16.0` tag was used when they were cloned and registered during the
`clone` command (i.e. `-r -t v0.16.0`)

```
> ziti-git unregister-tag v0.16.0
> rm -rf ./ziti-0.16.0/*
```

## Prior Art

Ziti Git is based off of [gmg](https://github.com/abrochard/go-many-git)
which in turn was inspired by the amazing
[mr](https://myrepos.branchable.com) and
[gr](https://github.com/mixu/gr) tools.

A big thanks to all.
