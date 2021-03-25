import gc
import esp
import network
import microcoapy
import ujson

WIFI_SSID = "Borg"
WIFI_PASSWORD = "startrekborg"

COAP_SERVER_ADDRESS = "192.168.64.105"
COAP_SERVER_PORT = 5683

esp.osdebug(None)
gc.collect()

# Connect to Wifi
wlan = network.WLAN(network.STA_IF)
wlan.active(True)
wlan.connect(WIFI_SSID, WIFI_PASSWORD)
wlan.ifconfig()


# CoAP callback
def receivedMessageCallback(packet, sender):
    print('Message received:', packet.toString(), ', from: ', sender)


# CoAP send message
def sendPostRequest(client):
    file = open("message.json", "r")
    message = file.read()
    file.close()

    messageId = client.post(
        COAP_SERVER_ADDRESS,
        COAP_SERVER_PORT,
        "/message",
        message,
        None,
        microcoapy.COAP_CONTENT_FORMAT.COAP_APPLICATION_JSON
    )

    print("[POST] Message Id: ", messageId)
    client.poll(10000)


# Intialize CoAP client
client = microcoapy.Coap()
client.discardRetransmissions = True
client.debug = True
client.responseCallback = receivedMessageCallback

client.start()

print("1   req.start")
sendPostRequest(client)
print("2   req.end")

client.stop()
