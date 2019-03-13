# git-semver: A cli tool for versioning your git repository

## Why

There is no standardized way to version a git repository. Node applications store their version in a package.json, other appications store it in a Makefile or somewhere else.

What i want is small cli tool that handles the version of my repository, no matter what technology is under the hood.

## Installation

**brew**

```bash
brew tap meinto/git-semver https://github.com/meinto/git-semver
brew install meinto/git-semver/git-semver
```

**manually**

Download the corresponding [latest binary](https://github.com/meinto/git-semver/releases) and run the `install` command. Right now the `install` command is only valid for mac and linux.

```bash
<name-of-binary> install
```

You can use `semver` as a gitplugin by using the following syntax:

```bash
git semver version ...
```

## Usage

By default this cli tool uses a `VERSION` file in the root folder of a git repository to store the version. The versioning is built up on [semantic versioning](https://semver.org/).

```bash
semver version [major|minor|patch] \
  [--dryrun] \                        # default: false -- only show how version would change
  [-p <path-to-repo>] \               # default: .
  [-f <version-file-name>] \          # default: VERSION -- define alternative version file
  [-t <version-file-type>] \          # default: raw -- you can set the values "json" or "raw"
  [--tag] \                           # default: false -- tag the commit with the new version
  [--push] \                          # default: false -- push all changes made by semver
```

You can create an individual `semver.config.json` file in the root of your project to override the default values of the flags. Simply run the `semver init` command and follow the instructions.

### Custom version file

By default `semver` lookup the `VERSION` file in the current directory. If you want to store your version in a custom file use the flag `-f`.

If you prefer a `json` file which contains the version number in a `version` property you can do this by using the flags `-f` in combination with `-t`:

> Auto file type detection will be implemented in the next release.  
> If you can't wait, please send a pull request :)

```bash
semver version minor -f package.json -t json
```

If your defined version file is of type `json`, the cli tool only overrides the property `version` and leaves other properties untouched.

### Create git tag

With the flag `--tag` or short `-T`, `semver` will create a git tag of the new version e.g.: `v1.0.0`.

### Push version changes

`semver` writes the new version back into the version file. As described, you can also tag the commit using the flag `-T`. To automatically push these changes made by `semver`, use the flag `--push` or short `-P`.

## Get Version(s)

With the `get` command you can get the current or next possible versions.

```bash
semver get         # will print the current version
semver get major   # will print the next major version
semver get minor   # will print the next minor version
semver get patch   # will print the next patch version
```

With the `--raw` or short `-r` flag you will get the plain number without description. For example:

```bash
semver get minor --raw
```

## Contribute

Create a fork, make your changes, and send a pull request. :sunglasses: