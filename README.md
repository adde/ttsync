# Toggl Tidig Sync

_Toggl Tidig Sync_ is a CLI tool that fetches time entries from __Toggl__(REST API) and writes them to a CSV file in a format that can be imported in __Tidig__ without having to do any manual adjustments.

## Prerequisites

You need to have the exact same namings of clients and projects in Toggl as in Tidig for this import to work seamlessly.

Also if you want to get correct activities in the import, you need to use tags in Toggl that matches activity names in Tidig.

To make use of the ticket number field, you can separate your descriptions in Toggl with a pipe(|). The part before the pipe will be used as the ticket number and the part after will be used as the description. Example description with a ticket number:
```
PROJ-123 | Did some work
```

## Installation

Make sure you have _go_ installed on your machine. Verify by running the following command in a terminal:
```sh
go version
```
It should return the installed version of _go_, otherwise you need to install go: https://go.dev/doc/install

To build executables for the app, you also need to have _make_ installed. Verify with this command:
```sh
which make
```
It should return something like this: `/usr/bin/make`. Otherwise you need to install make:
```sh
# MacOS
brew install make

# Ubuntu / WSL
sudo apt install make
```

## Usage

### Environment variables

#### Alternative 1 (development)

If you are only running the app from source, it is enough to just keep a `.env` file in your working directory. Create a copy of `.env.example` and save it to `.env`.

Replace `USERNAME` and `PASSWORD` placeholders with your own Toggl credentials.

#### Alternative 2 (production)

If you are planning on running the app from an executable, then you may want to create a permanent `.env` file in your home config directory:

Create a copy of `.env.example` and save it to `~/.config/ttsync/.env`.
```sh
cp .env.example ~/.config/ttsync/.env
```
Replace `USERNAME` and `PASSWORD` placeholders with your own Toggl credentials.

### Run

To run the app from source, run the following command in the root of the project:
```sh
go run src
```

### Build

To build executables from source, run the following command:
```sh
make build
```
This will create executables for all platforms in the `bin` directory.

Then move the correct binary for your platform to a folder in your PATH, for example:
```sh
mv ./bin/ttsync-linux-amd64 /usr/local/bin/ttsync
```

and then run the app:
```sh
ttsync
```

## Arguments

The following flags can be used when running the app to alter the output.

### start

The start date is used to filter which time entries to fetch from the Toggl API. (Default: previous monday)
```sh
ttsync -start 2023-11-01
ttsync -s 2023-11-01
```

### end

The end date is used to filter which time entries to fetch from the Toggl API. (Default: 2100-01-01)
```sh
ttsync -end 2023-11-03
ttsync -e 2023-11-03
```

### output

The output file path is used to write the CSV file. (Default: time-entries.csv)

```sh
ttsync -output entries.csv
ttsync -o entries.csv
```

## Download

If you don't want to run or build the app from source, you can download the latest version of the app from the [release section](https://github.com/adde/toggl-tidig-sync/releases) in this repo.
