## Release process

1. Merge all pr's to master which need to be part of the new release
2. Push a tag following semantic versioning prefixed by 'v'. Do not create a github release, this is done automatically.
3. At this point the release is done and the artifacts are getting built. The helm chart is automatically released, no further action needed. The kustomize base will automatically receive a renovate bump pr.