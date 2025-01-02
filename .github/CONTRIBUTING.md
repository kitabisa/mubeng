# Contributing

By participating to this project, you agree to abide our [code of conduct](https://github.com/mubeng/mubeng/blob/master/.github/CODE_OF_CONDUCT.md).

## Development

For small things like fixing typos in documentation, you can [make edits through GitHub](https://help.github.com/articles/editing-files-in-another-user-s-repository/), which will handle forking and making a pull request (PR) for you. For anything bigger or more complex, you'll probably want to set up a development environment on your machine, a quick procedure for which is as folows:


### Setup your machine

`mubeng` is written in [Go](https://golang.org/).

Prerequisites:

- make
- [Go 1.15+](https://golang.org/doc/install)

Fork and clone **[mubeng](https://github.com/mubeng/mubeng)** repository.

A good way of making sure everything is all right is running the following:

```bash
▶ make build
▶ ./bin/mubeng -h
```

### Test your change

When you are satisfied with the changes, we suggest you run:

```bash
▶ make test-extra
```

Which runs all the linters and tests.

### Submit a pull request

As you are ready with your code contribution, push your branch to your `mubeng` fork and open a pull request against the **master** branch.

Please also update the [CHANGELOG.md](https://github.com/mubeng/mubeng/blob/master/CHANGELOG.md) to note what you've added or fixed.

### Pull request checks

By submitting your pull request, you are agreeing to our [Contributing License Agreement](https://github.com/mubeng/mubeng/blob/master/.github/CONTRIBUTION_LICENSE_AGREEMENT.md).

Also, we run a few checks in CI by using GitHub actions, you can see them [here](https://github.com/mubeng/mubeng/tree/master/.github/workflows).