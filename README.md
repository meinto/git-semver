# git-semver: A cli tool for versioning your git repository

## Why

There is no standardized way to version a git repository. Node applications store their version in a package.json, other appications store it in a Makefile or somewhere else.

What i want is small cli tool that handles the version of my repository, no matter what technology is under the hood.

## Installation

Download the corresponding [latest binary](https://github.com/meinto/git-semver/releases) and run the `init` command. Right now the `init` command is only valid for mac and linux.

```bash
<name-of-binary> init

# answer the questions
Name of binary: <name-of-binary>
How do you want to use semver?
  > global
  > git plugin
```

You can use `semver` as a gitplugin by using the following syntax:

```bash
git semver version ...
```

## Usage

By default this cli tool uses a `semver.json` in the root folder of a git repository to store the version. The versioning is built up on [semantic versioning](https://semver.org/).

```bash
semver version [major|minor|patch] \
  [--dryrun] \                        # default: false -- only show how version would change
  [-p <path-to-repo>] \               # default: .
  [-o <name-of-version-json-file>] \  # default: semver.json -- define alternative version json file
  [--tag] \                           # default: false -- tag the commit with the new version
  [--push] \                          # default: false -- push all changes made by semver
  [-a <name-of-author>] \             # default: semver -- (only relevant when --push is set)
  [-e <email-of-author]            # default: semver@no-reply.git -- (only relevant when --push is set)
```

### Custom version file

By default `semver` lookup the `semver.json` in the current directory. If you want to store your version in a custom json file use the flag `-o`. The cli tool only overrides the property `version` and leaves other properties untouched.

### Create git tag

With the flag `-t`, `semver` will create a git tag of the new version e.g.: `v1.0.0`.

### Push version changes

`semver` writes the new version back into the config file. As described, you can also tag the commit using the flag `-t`. To automatically push these changes made by `semver`, use the flag `-P`.

The default author of the commit would be "semver" and the email "semver@no-reply.git". To change this, provide the flags `-a` (auhtor) and `-e` (email).

## Contribute

Create a fork, make your changes, and send a pull request. :sunglasses: