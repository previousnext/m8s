# package

Packages an app container image. The packaging steps are defined in the project `Dockerfile`.

<pre>
<code>{{ "skpr package --help 2>&1" | exec }}</code>
</pre>

## Examples

Tag the project and package using that tag.

```bash
$ git tag v1.0.0
$ git push --tags
$ skpr package v1.0.0
```

Package the current directory and tag with the current `git describe` value.

```bash
$ skpr package $(git describe --tags --always)
```

Package a project in another directory.

```bash
skpr package v1.0.0-alpha1 ~/code/client/project
```
