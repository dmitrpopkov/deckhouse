## Patches

### 001-remote-and-local-path-options.patch

Added 2 parameters for proxy operation mode:
- `remotepathonly`;
- `localpathalias`;

Example:
```yaml
proxy:
  remoteurl: "..."
  remotepathonly: "sys/deckhouse-oss"
  localpathalias: "system/deckhouse"
  username: "..."
  password: "..."
  ttl: 72h
```
Allows you to specify the allowed path to the registry, as well as replace the path for accessing the caching (local) registry

```bash
# not 'docker pull localhost:5001/sys/deckhouse-oss/install:latest'
docker pull localhost:5001/system/deckhouse/install:latest
```


### 002-scheduler-state-file-filling-and-deleting.patch

Adds logic for working with `/scheduler-state.json` file:
- For proxy mode, if the file is empty, a background job is started to fill it;
- If the mode is not proxy, the file is deleted;

It is necessary to switch from `Detached` to `Proxy` registry mode.
