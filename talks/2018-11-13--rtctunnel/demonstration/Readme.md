## [1] Start Client

```bash
docker run -it rtctunnel-demonstration
rtctunnel init && rtctunnel info
```

## [2] Start Server

```bash
docker run -it rtctunnel-demonstration
rtctunnel init && rtctunnel info
```

## [-] Add Routes

```bash
export CLIENT_KEY=DkXsGkKpZgwsNVE1QPYKmH8aRkj1NNK795pRMfqa82xK
export SERVER_KEY=GFxSiSSY5eHLMQ9sxSRA215Mf8WQ6MXxvokRqyW7WGEU
rtctunnel add-route \
    --local-peer=$CLIENT_KEY \
    --local-port=6379 \
    --remote-peer=$SERVER_KEY \
    --remote-port=6379
```

## [2] Start Redis Server

```bash
redis-server &
```

## [-] Start RTCTunnel (both)

```bash
rtctunnel run 2>/dev/null &
```

## [1] Start Redis Client

```bash
redis-cli INFO
```