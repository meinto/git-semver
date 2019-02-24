# git-semver: A cli for versioning your git repository

## Why

There is no standardized way to version a git repository. Node applications store their version in a package.json, other appications store it in a Makefile or somewhere else.

What i want is small cli that handles the version of my repository, no matter what technoligy is under the hood.

## Usage

By default this cli uses a `semver.json` in the root folder of a git repository to store the version. The versioning is built up on [semantic versioning](https://semver.org/).

```bash
semver version [major|minor|patch] \
  [--dryrun] \                        # default: false -- only show how version would change
  [-p <path-to-repo>] \               # default: .
  [-o <name-of-version-json-file>] \  # default: semver.json -- define alternative version json file
  [--tag] \                           # default: false -- tag the commit with the new version
  [--push] \                          # default: false -- push all changes made by semver
  [-a <name-of-author>] \             # default: semver -- (only relevant when --push is set)
  [-e <email-of-author]               # default: semver@no-reply.git -- (only relevant when --push is set)
```

### Custom version file

By default `semver` lookup the `semver.json` in the current directory. If you want to store your version in a custom json file use the flag `-o`. The cli only overrides the property `version` and leaves other properties untouched.

### Create git tag

With the flag `-t`, `semver` will create a git tag of the new version e.g.: `v1.0.0`.

### Push version changes

`semver` writes the new version back into the config file. Like described it also can tag the commit using the flag `-t`. To automatically push these changes made by `semver` use the flag `-P`.

The default author of the commit would be "semver" and the email "semver@no-reply.git". If you want to change this, provide the flags `-a` (auhtor) and `-e` (email).