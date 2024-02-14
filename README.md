# Easy-Transfer

## About
Simple ad-hoc file server for retrieving files over local network
via web upload that's accessible from all kinds of devices.

Goals:
- ease of use: just double-click the .exe and browse the website on the device
- transfer with full network speed

### Usage

**Info: You need Go installed on your system.** Ready-to-use builds will be provided soon. 
See [below](#install-go-tools).


Clone the project and build it manually:

```sh
git clone https://github.com/chucnorrisful/easy-transfer.git
cd easy-transfer
go build

# defaulting to target-folder "data"
easy-transfer.exe [<target-folder>]
```

Then a server will spin up to receive your uploaded files
and write them to the newly created directory.

**Info: How to open a cmd.exe console in Windows**
You need to type the command into a cmd.exe console.
Press Windows-key + R, the type cmd.exe and accept with Enter.
The console window should open immediately afterwards.

### Install Go Tools
Install the Go developer tools on your PC if you haven't done already.

**For Linux/Debian**
```sh
sudo apt-get update && sudo apt install golang-go 
```

**For Windows**

Browse the official download page of [Go](https://go.dev/dl/) and install a proper
version of the Go developer tools for your PC.