
## Sending messages

### libcoap2 (prefered)

- has TLS support

```sh
sudo apt install libcoap2-bin
coap-client -m put coap://[::1]/message -f message.json
```

### coap-cli

```sh
# https://github.com/avency/coap-cli
npm install coap-cli -g
```

```sh
echo -n '{ "deviceId": "7a4dd0ca-5ef6-490e-8373-1bf5d428f24a","data": [{"metric": "temperature","value": 22.3,"timestamp": 1616637069},{"metric": "pressure","value": 1000,"timestamp": 1616637069}]}' | coap post coap://localhost/message

cat message.json | coap post coap://localhost/messag
```