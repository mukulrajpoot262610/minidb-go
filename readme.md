# MiniDB

MiniDB is a lightweight, file-based relational database implemented in Go. It provides a simple command-line interface for creating tables, inserting, selecting, and deleting data, and managing schemas, all stored in local files. This project is ideal for learning about database internals, file I/O, and Go programming.

## Features

- Create tables with custom schemas (INT and TEXT columns)
- Insert, select, and delete rows
- Show and drop tables
- Persistent storage using binary files for data and JSON for schemas
- Simple SQL-like command interface

## Getting Started

### Prerequisites

- Go 1.24 or later

### Installation

Clone the repository and build the project:

```sh
git clone <repo-url>
cd basics
go build -o minidb main.go
```

### Usage

Run the database CLI:

```sh
./minidb
```

You will see:

```
MiniDB v0.1
Type 'help' for available commands.
Type 'exit' to exit the program.
db >
```

## Supported Commands

- `create table <table_name> (<column_name> <column_type>(<size>))`
  Example: `create table users (id INT, name TEXT(32), email TEXT(255))`
- `insert into <table_name> values (<values>)`
  Example: `insert into users values (1, "Alice", "alice@email.com")`
- `select * from <table_name>`
- `delete from <table_name> where <id>`
- `show tables`
- `drop table <table_name>`
- `help` — Show all commands
- `exit` — Exit the program

## Project Structure

```
basics/
  commands/      # Command implementations (select, insert, delete, etc.)
  constants/     # Schema and type definitions
  core/          # Command parser and dispatcher
  db/            # Data files (created at runtime)
  schemas/       # Table schemas (JSON, created at runtime)
  utils/         # Utility functions (file I/O, serialization, etc.)
  main.go        # Entry point
  go.mod         # Go module file
  readme.md      # This file
```

## How It Works

- **Schemas** are defined and stored in `schemas/schema.json`.
- **Data** for each table is stored in a binary file in `db/<table>.db`.
- On startup, schemas are loaded from file. Data is loaded on demand.
- All commands are parsed and dispatched via the REPL in `main.go` and `core/parser.go`.

## Development Notes

- Only `INT` and `TEXT(size)` column types are supported.
- All data is stored locally; there is no network or concurrency support.
- The code is modular for easy extension (add new commands, types, etc.).

## License

MIT
