# Toggl Tidig Sync

Fetch time entries from Toggl(API) and write the data to a CSV file that can be imported in Tidig without any need of manual adjustments.

## Prerequisites

You need to have the exact same namings of clients and projects in Toggl as in Tidig for this import to work seamlessly.

Also if you want to get correct activities in the import, you need to use tags in Toggl with the same activity names as in Tidig.

To make use of the ticket number field, you can separate your descriptions in Toggl with a pipe(|). The part before the pipe will be used as the ticket number and the part after will be used as the description.

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
go run ./src
```

## Build

To build an executable from source, run the following command:
```sh
go build ./src/main.go
```

Then move the binary to a folder in your path, for example:
```sh
mv main /usr/local/bin/ttsync
```
and then run the app:
```sh
ttsync
```
(you will still need the `.env` file in your working directory)