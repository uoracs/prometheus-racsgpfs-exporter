# prometheus-racsgpfs-exporter

This is a quick-and-dirty exporter so we can graph GPFS Fileset metrics such as
Size (GB), Quota (GB), and Inodes.

## Installation

Clone this repository

`git clone https://github.com/uoracs/prometheus-racsgpfs-exporter`

Navigate into the cloned directory

`cd prometheus-racsgpfs-exporter`

Build and install the application

`make install`

## Configuration

Available environment variables for configuration

#### `RACSGPFS_EXPORTER_LISTEN_ADDRESS`

Set a specific listen address for the metrics server.

default: `:8030`

## Metrics

#### `racsgpfs_size_gb`

Current size in GB of the fileset

#### `racsgpfs_quota_gb`

Current quota in GB of the fileset

#### `racsgpfs_inodes`

Current inode count of the fileset

## Uninstall

`make clean`
