# udmp-wpa-supplicant-monitor

Monitors the wpa-supplicant container status and starts it if not running.

## Build

```
docker build -t udmp-wpa-supplicant-monitor .
```

## Executing the container

```
docker run --rm -it -v "$HOME/pwd-file.txt":/pwd-file.txt \
    udmp-wpa-supplicant-monitor \
    root 192.168.1.1:22 /pwd-file.txt
```

## Executing as Corn task

```
0 */8 * * * docker run --rm -v "$HOME/pwd-file.txt":/pwd-file.txt ashrithr/udmp-wpa-supplicant-monitor:v0.0.1 root 192.168.1.1:22 /pwd-file.txt
```