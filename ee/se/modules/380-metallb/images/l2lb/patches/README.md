# Patches

## 001-l2lb-speaker-preffered-node.patch

Added preferred L2 speaker node feature

Upstream <https://github.com/metallb/metallb/pull/2246/>

## 002-l2lb-service-custom-resource.patch

Added a custom resource L2LBService to replace the original Service.

The controllers logic is rewritten to watch this new private resource.
