# Project: electronic-health
This repository contains an electronic health application that provides a user interface for creating, modifying, and viewing patient records.

## Notice
This application is currently under development and should not be used as a tool in actual healthcare information systems. The project is a part of 100 Days of Code challenge being undertaken by the software developer Marcie Henderson with the purpose of demonstrating/improving their professional skillset.

## Usage

### Install Source

- Execute the following command to clone repository:

```bash
git clone https://github.com/marciehenderson/electronic-health.git
```

### Database

- Install MySQL and create a password enabled user

- Log into MySQL from the database directory and enter the following to create the database and enter it automatically:

```sql
source ./generate-db.sql
```

- After the database has been successfully generated, test data may automatically populate the database tables by executing `go run .` while within the database directory.

### Server

- The server requires GoLang version *go1.22.5 linux/amd64*, but may work with other recent Go versions *(untested)*.

- Before running the server remember to export the required database credentials to the following environment variables *(replace username & password with credentials)*:

```bash
export DBUSER='username' && export DBPASS='password'
```

- Starting the server is as simple as executing the command `go run .` while inside the backend directory. *(Note: must compile frontend prior to first run)*

### Client

- Compiling `app.ts` to `bundle.js` is accomplished with the following command while in the frontend directory *(requires npm, npx, tsc, and rollup)*:

```bash
npx tsc && npx rollup -p @rollup/plugin-node-resolve ./dist/app.js -o ./dist/bundle.js
```

- The user interface is accessible at the url which appears in the server console on runtime. *(i.e., `https://localhost:8080`)*