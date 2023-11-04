# Toggl Tidig Sync

_Toggl Tidig Sync_ is a CLI tool that fetches time entries from __Toggl__(REST API) and writes them to a CSV file in a format that can be imported in __Tidig__ without having to do any manual adjustments.

## Prerequisites

You need to have the exact same namings of clients and projects in Toggl as in Tidig for this import to work seamlessly.

Also if you want to get correct activities in the import, you need to use tags in Toggl that matches activity names in Tidig.

To make use of the ticket number field, you can separate your descriptions in Toggl with a pipe(|). The part before the pipe will be used as the ticket number and the part after will be used as the description. Example desciption with a ticket number:
```
PROJ-123 | Did some work
```

## Installation

Make sure you have go installed on your machine. Verify by running the following command:
```sh
go version
```
It should return the installed version of go, otherwise you need to install go: https://go.dev/doc/install

## Usage

Create a copy of `.env.example` and name it `.env`.

Replace `USERNAME` and `PASSWORD` placeholders with your own Toggl credentials. You can also adjust the `FROM` value to the date where you want to start fetch entries.

To run the app from source, run the following command in the root of the project:
```sh
go run src
```

## Build

To build executables from source, run the following command:
```sh
make build
```
This will create executables for all platforms in the `bin` directory.

Then move the binary to a folder in your path, for example:
```sh
mv ./bin/ttsync-linux-amd64 /usr/local/bin/ttsync
```
and then run the app:
```sh
ttsync
```
(you will still need the `.env` file in your working directory when running the app)

## Download

Alternatively, you can download the latest version of the app from the [release section](https://github.com/adde/toggl-tidig-sync/releases) in this repo.

## Changelog

### v1.0.0 (2023-11-03)

* First release.