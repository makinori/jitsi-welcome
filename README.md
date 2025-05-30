# Jitsi Welcome

Welcome page for my Jitsi instance

https://jitsi.hotmilk.space

Grabs all the character names from all the anime I've seen, caches them and uses them for generating a random room name.

`DEV=1 go run .`

### Docker Compose and Traefik

Example service:

```yml
welcome:
    build: ./jitsi-welcome
    restart: ${RESTART_POLICY:-unless-stopped}
    environment:
        PORT: 8080
        ANILIST_USERNAME: makinori
        CACHE_PATH: /cache/cache.json
    volumes:
        - ./jitsi-welcome-cache:/cache
    labels:
        - service=jitsi-web
        - traefik.enable=true
        - >
            traefik.http.routers.jitsi-welcome.rule=
            Host(`jitsi.hotmilk.space`) &&
            (Path(`/`) || PathPrefix(`/welcome/`))
        - traefik.http.routers.jitsi-welcome.entrypoints=websecure
        - traefik.http.routers.jitsi-welcome.service=jitsi-welcome
        - traefik.http.services.jitsi-welcome.loadbalancer.server.port=8080
        - traefik.http.routers.jitsi-welcome.tls.certresolver=le
        - traefik.docker.network=traefik
    networks:
        - traefik
```
