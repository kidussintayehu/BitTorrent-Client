
 # BitTorrent Client Implementation in Go

This is a BitTorrent client implementation in Go language. The client is capable of downloading and uploading files using the BitTorrent protocol.BitTorrent is a protocol for downloading and distributing files across the Internet. It is a peer-to-peer protocol because users of the BitTorrent network, known as peers, download data in parts from one another. Peers are introduced to one another by trackers, which are central servers.

Seeder is a user that is seeding the file which means they allow other users to take bits of their data that has been fully downloaded in their computer. 

## Requirements

- Go version 1.16 or later
- Git

## Installation

1. Clone this repository:

 
git clone https://github.com/kidussintayehu/BitTorrent-Client.git
 

2. Change to the directory where the repository was cloned:

 
cd BitTorrent-Client
 

3. Build the client using the  go build  command:

 
go build
 

## Usage

1. Run the client with the following command:

go run main.go <path-to-torrent-file> <path-to-download-directory>
 

For example:
 
./bittorrent-client ~/testfiles/debian-10.2.0-amd64-netinst.iso.torrent  ~/downloads/
 

2. The client will start downloading the file. The progress of the download will be displayed in the terminal.

3. Once the download is complete, the file will be saved to the download directory specified in the command.

## Features

- Download and upload files using the BitTorrent protocol.
- Support for multiple peers.
- Ability to resume downloads.


