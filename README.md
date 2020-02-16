# Trinity
Usefull framework to operating Web Applications for beginners.

## What's this?
This Application is framework to solve trouble in operating Web services.

Individual developers doesn't have
- Money
- Time
- manpower

Trinity has
- Logging
- Notification using Webhook
- Continuous Delivery with GitHub

Trinity reduce your workload.

### Prerequisites
You need to Build trinity
- Node.js (above v12.4)
- Go (above v1.11 )
    - [rakyll/tatik](https://github.com/rakyll/statik)


## Requirements and Setup
You can try Trinity by downloading from [GitHub - Releases]().

Or You can build yourself.
### Cloning the Repo
Before you start working on Trinity, you'll need to clone our GitHub repository:

```sh
git clone git@github.com:sechack-z/trinity
```

Now, enter the repository.

```sh
cd trinity
```
### Build as single binary
We want to make easy to use Trinity.
So Trinity can build as single binary.
#### Install Dependencies.
- Node.js (above v12.4)
- Go (above v1.11)
- [rakyll/tatik](https://github.com/rakyll/statik)
- make

or server and client can be built separately.

### Client 
Client works in `client` directory.
```
cd client
```
#### Node.js
Client is written by [Node.js](https://nodejs.org/en/).
You'll need to install `Node.js` for your system to build client.

**Note**: This is confirmed to work on `Node.js v12.4`, but there may be issues on other versions. If you have trouble, please bump your Node.js version to 12.4.
#### Installing Dependencies
We're using a package manager called Yarn. You'll need to install Yarn before continuing.


Install all required packages with:

```sh
yarn
```

#### Building Client
Build client:

```sh
yarn build:stage
```


### Server
Server works in `server` directry.
#### Go
Server is written by [Go](https://golang.org).
You'll need to install `Go` for your system to build server.

**Note**: This is confirmed to work on `Go v1.11`, but there may be issues on other versions. If you have trouble, please bump your Go version to 1.11.
#### Installing Dependencies
This project is using `Go modules`.
So dependencies is automatically download when you build server.

#### Building Server
```sh
go build
```

