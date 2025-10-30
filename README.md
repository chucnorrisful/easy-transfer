# Easy-Transfer

## About
This project aims to solve one simple problem: To quickly send files
from your phone to your computer! And to do so:
- device & OS independent
- with maximum transmission rate
- very user friendly

## Usage (Windows)
After installation (see [below](#install)), the app is ready for use:

- Double-click **easy-transfer.exe** on your PC
- There might be a windows security pop-up to trust the .exe
- There might be a firewall pop-up to allow acces to your network - that is required for the file-transfer to work
 <br>-> a window with a QR code will appear
- Scan the QR code with your phone and open the link
 <br>-> an upload website will open on your phone
- Select files / folders to upload
- Press the colorful "Upload" Button
 <br>  -> the selected files will be copied to your PC
- Find uploaded files on your PC in the ./data directory (next to **easy-transfer.exe** file)

*Linux & MacOS: the executable does not have a .exe file ending, otherwise it is very similar.*

## Capabilities / Limitations
The transmission only works if your phone / device is in the same
local network (WLAN, LAN) as your PC.
In case your sending device is a PC, this also works for PC to PC transmissions.
Open the link of the upload website on the sending PC via browser.

## Technical Details
The .exe serves a website to be accessed by the sending devices
plus a simple file server for retrieving files over local network.

The sender website uses standard, multi-platform browser APIs for upload.
To keep the app lightweight and easy to maintain, it's written in Go + vanilla 
Javascript + Axios for the file transmission.

### Install
You can build the program from source, or download the .exe from the latest
[release](https://github.com/chucnorrisful/easy-transfer/releases/latest).

**Info: For building from source, you need Go installed on your system.**
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

Browse the official download page of [Go](https://go.dev/dl/) and install the Go compiler for your PC.

**For Linux/Debian**

```sh
sudo apt-get update && sudo apt install golang-go
```
