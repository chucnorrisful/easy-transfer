# Easy-Transfer

## About
Simple ad-hoc file server for retrieving files over local network
via web upload that's accessible from all kinds of devices.

Goals:
- ease of use: just double-click the .exe and browse the website on the sending device
- transfer with full network speed

### Usage

**Info: You need Go installed on your system.** Ready-to-use builds will be provided soon. 
See [below](#install-go-tools).

Clone the project and build it manually:

```sh
git clone https://github.com/chucnorrisful/easy-transfer.git
cd easy-transfer
go build

# run the executable (or just double-click it in the explorer)
./easy-transfer.exe
```

Then a server will spin up to receive your files
and write them to a newly created directory /data.

### Install Go

Install the Go programming language on your PC.

**For Linux/Debian**
```sh
sudo apt-get update && sudo apt install golang-go 
```

**For Windows**

Browse the official download page of [Go](https://go.dev/dl/) and install a proper
version of the Go developer tools for your PC.