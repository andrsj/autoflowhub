# AutoFlowHub

AutoFlowHub is a powerful tool designed to automate workflows, track GitHub releases, deploy devnet instances, prepare releases on GitHub, and push them, format, and edit files, and much more.

## Features

- Automate workflows and track GitHub releases
- Deploy devnet instances
- Prepare and push releases on GitHub
- Format and edit files

## Project Structure

```
autoflowhub/
├── cmd
│ ├── change_notifier
│ │ └── main.go
│ └── release_fetcher
│ └── main.go
├── docker-compose.yml
├── internal
│ ├── adapters
│ │ ├── adapters_test.go
│ │ └── github.go
│ ├── database
│ │ └── storage.go
│ ├── models
│ │ └── release.go
│ ├── notifications
│ │ ├── email.go
│ │ ├── notifications_test.go
│ │ ├── slack.go
│ │ └── sms.go
│ └── utils
│ ├── authentication.go
│ └── logging.go
├── pkg
│ ├── api
│ │ ├── change_notifier
│ │ │ └── change_notifier.go
│ │ └── release_fetcher
│ │ └── release_fetcher.go
│ └── config
│ └── config.go
└── README.md
```