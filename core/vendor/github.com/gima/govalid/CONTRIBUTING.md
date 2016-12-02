## Contributing

File an issue or use pull requests, but relay your intent for the action.

Try to avoid mixing different concerns in one commit. *(Look who's talking.)* Same applies to issues and pull requests.

*By contributing to this project, you agree to place (and you agree that you have the right to place) your contribution under the same license as this project is.*



#### Code formatting

Run `go fmt` before committing.

#### Run tests

Test all reasonable code paths.

    go test ./tests/

#### Run code coverage

Try to cover the necessary cases.

    go test -cover -coverpkg . -coverprofile cover.prof ./tests
    go tool cover -html=cover.prof -o coverage.html
    
---

*This contribution readme was shamelessly modelled after:  
https://opencomparison.readthedocs.org/en/latest/contributing.html*
