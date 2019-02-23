# git-semver: A tool for versioning your git repository

## Why

There is no standardized way to version a git repository. Node applications store their version in a package.json, other appications store it in a Makefile or somewhere else.

What i want is small cli that handles the version of my repository, no matter what technoligy is under the hood.

## Usage

By default this cli uses a `semver.json` in the root folder of a git repository to store the version. The versioning is built up on [semantic versioning](https://semver.org/).

```bash
semver version [major|minor|patch] [--dryrun] [-p <path-to-repo>] [-o <name-of-version-json-file>]
```

By default `semver` lookup the `semver.json` in the current directory. If you want to store your version in a custom json file use the flag `-o`. The programm only overrides the property `version` and leaves other properties untouched.