# gh-cleaner

An interactive Go CLI for reviewing and deleting your GitHub repositories. A Bayesian classifier uses locally recorded decisions to order the repositories for review; you confirm each deletion yourself.

**Status:** experimental. Start with `--dry-run` to learn the workflow without deleting repositories.

![Interactive repository review](assets/gh-cleaner-show.gif)

## Build

Requires Git and Go 1.22.2 or newer.

```sh
git clone https://github.com/KitsuneSemCalda/gh-cleaner.git
cd gh-cleaner
go build -o gh-cleaner ./cmd/gh-cleaner
./gh-cleaner --help
```

With Make installed, `make install` builds and copies the executable to `~/.local/bin`. Add that directory to your `PATH` to run `gh-cleaner` from elsewhere.

## Authentication

The program reads your GitHub username and personal access token from `~/.netrc`. The current parser requires the GitHub entry on **one line**:

```text
machine github.com login YOUR_GITHUB_USERNAME password YOUR_GITHUB_TOKEN
```

Edit an existing file to add or update that entry, and restrict access to it:

```sh
chmod 600 ~/.netrc
```

Keep this file outside the repository. The token must have access to the repositories you want to review; deleting a repository also requires permission to delete it. The program does not use the GitHub CLI's stored login.

## Review without deleting

```sh
./gh-cleaner --dry-run
```

The program lists repositories owned by the authenticated user, including private repositories accessible to the token. Forks are excluded from the review by default.

For each repository, it displays details and asks whether to delete it. Answering `y` opens a second confirmation. In dry-run mode, neither confirmation sends a deletion request.

**Dry-run still records local training decisions.** These are stored under `~/.local/share/gh-cleaner/repository/`, in `saved/` and `deleted/`. These files contain repository metadata, not backups of repository contents; a record in `deleted/` does not prove that a repository was deleted.

On a first run, a warning about missing training data can appear before the review starts.

## Delete selected repositories

> [!CAUTION]
> Without `--dry-run`, confirming a deletion sends a real GitHub repository deletion request. Back up anything you need before proceeding. The second confirmation initially selects **No**; you must deliberately move to **Yes** before pressing Enter.

```sh
./gh-cleaner
```

To include forks in either mode:

```sh
./gh-cleaner --dry-run --forks
./gh-cleaner --forks
```

| Flag | Effect |
|---|---|
| `--dry-run` | Review and record training decisions without deleting repositories. |
| `--forks` | Include forked repositories in the review. |
| `--help` | Show command-line options. |

## Development

```sh
go build ./...
go vet ./...
go test ./...
```

There are currently no automated test cases. `go test ./...` checks that the packages compile, but does not validate the deletion workflow.

The CLI uses [go-github](https://github.com/google/go-github), [bayesian](https://github.com/jbrukh/bayesian), and [promptui](https://github.com/manifoldco/promptui). Contributions are welcome.

## License

See [LICENSE](LICENSE).
