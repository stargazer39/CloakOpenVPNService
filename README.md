# CloakOpenVPNService

## Run Cloak + OpenVPN on your linux pc with reconnect on failure for reliability
Run like this :
```bash
$ go build .
$ sudo ./cloakopenvpnservice --cloak-host 188.166.240.175 --cloak-port 443
```
Make sure you have openvpn and ck-client in PATH
