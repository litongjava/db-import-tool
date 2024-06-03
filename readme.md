# db-improt-tool

## Overview

db-improt-tool is a Go application designed to read data from an Excel file and insert it into a specified database table. It supports both PostgreSQL and MySQL databases, and automatically handles inserting data without hardcoding table structures.

## Features

- Supports PostgreSQL and MySQL databases.
- Reads data from an Excel file.
- Dynamically generates SQL insert statements based on the table headers.
- Ignores empty fields during insertion.
- Allows disabling SSL for PostgreSQL connections.

## Prerequisites

- Go 1.16 or higher
- PostgreSQL or MySQL database

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/litongjava/db-improt-tool.git
   cd db-improt-tool
   ```

2. Install the dependencies:
   ```sh
   go get -u github.com/xuri/excelize/v2
   go get -u github.com/cloudwego/hertz/pkg/common/hlog
   go get -u github.com/go-sql-driver/mysql
   go get -u github.com/lib/pq
   ```

## Usage

1. Compile the program:
   ```sh
   go build -o db-improt-tool main.go
   ```

2. Run the program with the required parameters:
   ```sh
   ./db-improt-tool -dsn="your_database_dsn" -excel="path_to_your_excel_file.xlsx" -table="your_table_name"
   ```

### Example Command

For PostgreSQL:
```sh
./db-improt-tool -dsn="postgresql://username:password@127.0.0.1:15432/dbname?sslmode=false" -excel="path_to_your_excel_file.xlsx" -table="your_table_name"
```

For MySQL:
```sh
./db-improt-tool -dsn="your_mysql_dsn" -excel="path_to_your_excel_file.xlsx" -table="your_table_name"
```

## Parameters

- `-dsn`: The Data Source Name (DSN) for the database connection. The format depends on the database type.
- `-excel`: The file path to the Excel file containing the data to be inserted.
- `-table`: The name of the database table where the data will be inserted.

## Code Structure

- `main.go`: Main application logic including initialization, database connection, reading Excel data, and inserting data into the database.

## Error Handling

The application includes basic error handling and logs any errors encountered during execution. If there is an issue with the database connection or executing the query, the error will be logged, and the program will exit.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.