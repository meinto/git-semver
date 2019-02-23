# git-semver: A tool for versioning your git repository

## Why

There is no standardized way to version a git repository. Node applications store their version in a package.json, other appications store it in a Makefile or somewhere else.

What i want is small cli that handles the version of my repository, no matter what technoligy is under the hood.

## Usage

By default this cli uses a package.json in the root folder of a git repository to store the version. The versioning is built up on [semantic versioning](https://semver.org/).

```bash
semver version [major|minor|patch] [-p <path-to-repo>]
```

By default `semver` lookup the `package.json` in the current directory. The programm only overrides the property `version` and leaves other properties untouched.