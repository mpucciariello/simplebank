## SimpleBank

*simple_bank* is a small project built with Golang, Postgres, and Docker. It is based on a series of [courses](https://www.youtube.com/watch?v=rx6CPDK_5mU&list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE&pp=iAQB) and aims to create a minimal version of a bank, simulating the interaction between users, accounts, and their transfers.

### Setup local environment
_Setup in macOS with Homebrew._

#### migrate
run `brew install golang-migrate`
#### sqlc
run `brew install sqlc`
#### gomock
run `go install github.com/golang/mock/mockgen@v1.6.0`
#### linter
run `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
#### pre_commit
run `pip install pre-commit` `pre-commit install`
