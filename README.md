# WIP: docker bittorrent graph driver

## Design

This project is an effort to implement a transparently shared bittorrent-based
graph driver for docker.

Each node will maintain it's layer cache as a set of directories containing the
layer diff, each one calculated as a torrent file (stored as the SHA256 content
hash ID).

Locally, the graph driver will replicate the functionality of the `overlay2`
driver to expose layers to docker. 

Behind the scenes, the driver makes all its layers available seamlessly to other 
connected hosts as part of a bittorrent swarm, allowing rapid bandwidth 
efficient distribution of image layers.

At least, that's the idea :) Let's see if I can do it. Intended to resolve
https://github.com/docker/docker/issues/25997 for me.
