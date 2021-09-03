# GitVersion
A helper tool for discerning Semver versions from Git tags / commits / branch names.

## How to derive a version
```bash
gitversion derive
```
This command will derive a proper semantic version ([SemVer](https://semver.org/)) from the latest commit tags or the branch name.

### Example
#### Getting the version from the branch name:
```bash
> git branch --show-current
release/1.0.2-rc.1
```

Will result in the following output:
```bash
> gitversion derive --pretty
{
  "major": 1,
  "minor": 0,
  "patch": 2,
  "core": "1.0.2",
  "isDefault": false,
  "preRelease": "rc.1",
  "metadata": "",
  "semver": "1.0.2-rc.1",
  "sha": "01b5e515606340d91f420c75a7301210f7eeb840",
  "shortSha": "01b5e51",
  "branchName": "release/1.0.2-rc.1"
}
```

#### Getting the version from the current commits list of tags.
This will return the highest semver compliant version:
```bash
> git log -n 1 --tags --format="%D"
HEAD -> main, tag: 1.0.0, tag: latest_version, tag: 0.1.0-test
```

Will result inthe following output:
```bash
> gitversion derive --pretty
{
  "major": 1,
  "minor": 0,
  "patch": 0,
  "core": "1.0.0",
  "isDefault": false,
  "preRelease": "",
  "metadata": "",
  "semver": "1.0.0",
  "sha": "01b5e515606340d91f420c75a7301210f7eeb840",
  "shortSha": "01b5e51",
  "branchName": "main"
}
```

#### Default, if all else fails.
```bash
> git tags

```
Will result in:
```bash
> gitversion derive --pretty
{
  "major": 0,
  "minor": 1,
  "patch": 0,
  "core": "0.1.0",
  "isDefault": true,
  "preRelease": "",
  "metadata": "",
  "semver": "0.1.0",
  "sha": "01b5e515606340d91f420c75a7301210f7eeb840",
  "shortSha": "01b5e51",
  "branchName": "main"
}
```
