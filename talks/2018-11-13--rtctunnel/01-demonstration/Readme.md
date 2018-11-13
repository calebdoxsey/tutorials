## [A] Start Client

```bash
docker run -it rtctunnel-demonstration
rtctunnel init && rtctunnel info
```

## [B] Start Server

```bash
docker run -it rtctunnel-demonstration
rtctunnel init && rtctunnel info
```

## [+] Add Routes

```bash
export CLIENT_KEY=
export SERVER_KEY=
rtctunnel add-route \
    --local-peer=$CLIENT_KEY \
    --local-port=6379 \
    --remote-peer=$SERVER_KEY \
    --remote-port=6379
```

## [B] Start Redis Server

```bash
redis-server &
```

## [+] Start RTCTunnel (both)

```bash
rtctunnel run &
```

## [A] Start Redis Client

```bash
redis-cli INFO
```