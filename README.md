# Easy-Transfer

This project aims to solve one simple problem: To quickly send files from your phone to your computer!
And to do so:
- device agnostically
- with maximum speed
- without a hustle

*This project focusses on Windows atm, but should already work on linux/mac - just without some quality of life features.*

## How to use:
(see [below](#install) on how to install)

- Double-click **easy-transfer.exe** on PC
- A QR Code appears
- Scan it with phone, open link
- Upload files via the web interface
- Find uploaded files on PC in ./data directory (opened automatically)

Only works, if your phone/device is in the same network (WLAN, LAN) as PC.
Also works for PC to PC transfers - the link for the upload site is copied to your clipboard on startup of *easy-transfer.exe*, just send it to the second PC!

## Technical how:
Simple ad-hoc file server for retrieving files over local network via web upload that's accessible from all kinds of devices. 
Uses standard browser APIs for upload, and a locally hosted Go backend server for hosting the upload site and receiving the files.

### Install

You can build the program from source, or download the .exe from the latest [release](https://github.com/chucnorrisful/easy-transfer/releases/latest).

**Info: You need Go installed on your system.**
See [below](#install-go).

Clone the project and build it manually:

```sh
git clone https://github.com/chucnorrisful/easy-transfer.git
cd easy-transfer
go build

# run the executable (or just double-click it in the explorer)
./easy-transfer.exe
```

### Install Go

Install the Go programming language on your PC.

**For Windows**

Browse the official download page of [Go](https://go.dev/dl/) and install a proper
version of the Go developer tools for your PC.

**For Linux/Debian**
```sh
sudo apt-get update && sudo apt install golang-go 
```
