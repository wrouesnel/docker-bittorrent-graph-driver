[![Build Status](https://travis-ci.org/wrouesnel/docker-bittorrent-graph-driver.svg?branch=master)](https://travis-ci.org/wrouesnel/docker-bittorrent-graph-driver)

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

## How it works

Internally this driver implements a slightly more efficient VFS graph
driver (to keep the design initially simple). Read-only layers are
managed with hard links, and read-write layers are treated as full
copies.

Since there is a need for multiple graph-drivers to advertise their
torrent availability, this driver does depend on a distributed key-value
store as well. Each node publishes it's list of available layers, and
queries the DKV to find out if a layer is available in the swarm.

Each layer produces 2 torrent files:
* A single-file torrent containing the complete layer description,
  including the magnet links for the containing layers.
* A multi-file torrent containing the difference layer and metadata
  for the specific layer.

When a connected docker daemon requests a given layer, the driver first
checks if the layer is already available locally (and adds a reference
file to the ref directory of the layer).

If the layer is not available, it then queries the DKV to find out if
the layer is known to the swarm and acquire the magnet URI of the
layer description.

The driver then begins downloading the layers it does not have, and
finally constructs and returns the layer docker requested.

Any layers the driver holds are automatically published to the DKV and
it acts as a seeding host for them for other members of the swarm.

Layers which are no longer referenced by any connected daemon, and which
meet the minimum seeding level for a layer to the swarm, are deleted
from the host.

The seeding level of a layer is maintained within the DKV. This is too
allow the eventual GC of layers which are not used within the swarm
anywhere, by setting their seeding level to 0 which will cause them to
eventually be deleted.

## Missing functionality
Currently a new layer being pulled from an upstream repository can
cause a storm if many docker agents try to pull it simultaneously. This
will be resolved by a "back-off" functionality in the DKV, which will
give an initial layer pull a chance to seed into the cluster so load
can be reduced on the other agents.