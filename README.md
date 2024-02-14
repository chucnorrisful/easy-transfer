# Easy-Transfer

## About
Simple ad-hoc, file server for retrieving files over local network
via website upload that's accessible from all kinds of devices.

Goals:
- ease of use: just double-click the exe and browse the website on the device
- transfer with full network speed

### Usage
Download easy-transfer.exe and double-click it.
Then a server will spin up to receive your uploaded files
and write it to the "data" directory next to the exe file.

```sh
# defaulting to target-folder "data"
easy-transfer.exe [<target-folder>]
```

**Info: How to open a cmd.exe console in Windows**
You need to type the command into a cmd.exe console.
Press Windows-key + R, the type cmd.exe and accept with Enter.
The console window should open immediately afterwards.

## Build From Source Code
The project is simple. Just install the Go developer tools and build an exe file.
You can then run the exe file by double clicking and should be good to go. (pun not intended ^^)

### Install Go Tools
Install the Go developer tools on your PC if you haven't done already.

**For Linux/Debian**

```sh
sudo apt-get update && sudo apt install golang-go 
```

**For Windows**

Browse the official download page of [Go](https://go.dev/dl/) and install a proper
version of the Go developer tools for your PC.

### Build



```sh
go build
```
