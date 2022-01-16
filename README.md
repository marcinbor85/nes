# NES Messenger
NES Messenger is an educational project (at least for the original author)
that aims to be the most secure instant messaging internet communicator as possible.
At the same time it has no need to use any advanced cloud-base services,
nor storing any private data. An uncompromising approach makes it unique.
It relies on the maximum extent on free, public and open source solutions.

It doesn't make sense to fight for a little bit of security where it's not absolutely necessary.
Therefore the messenger uses public MQTT brokers and asymmetric cryptography with publicly available keys.
Let what is public remain public, and what is private remain private, without exceptions.

## How it works?
NES Messanger uses two cloud-base services.

First one is public keys provider which as the name suggests, is intented to share public keys of registered users.
As You probably know, public keys are used to encrypt messages in asymmetric encryption algorithms like RSA-2048,
and for signature verification. So this is the most secure solution, as long as private key remains on the client's device.

The second one is a MQTT broker. Due to the high level of security and strong encryption,
it is safe to use even a public broker. Of course, anyone will be able to intercept encrypted messages
(which is generally obvious, even in the most secure environment), but only the original recipient of the message
will be able to decrypt it. Darkest is under the lantern. You could use Your own private broker, or use already existing public one.

![alt text](assets/cloud.png?raw=true "Cloud architecture")

The AES-256 algorithm with symmetric keys is used to encrypt messages. Symmetric keys are randomly generated
for each message and encrypted with the RSA algorithm with recipient public key. Therefore, without the private key,
it is not possible to restore the encrypted messages string, even if the AES key is broken by brute force.
Such an operation would have to be repeated for each subsequent message, which is an extremely complex
and time-consuming operation. For example, even if the interceptor manages to decrypt the first message
that you are dating someone for a beer, in order to find out what time it will be,
he will have to decrypt the second message encrypted with diffenet keys again.

In addition, each message is signed by sender private key, so the recipient can verify that the sender
is who they say they are. The signing and verification algorithm is triggered automatically for each message.
This makes the NES communicator resistant to Man In The Middle attacks as well.

## Flow
What does the message exchange between NES messenger clients look like?\
This is shown in the diagram below.

![alt text](assets/message_flow.png?raw=true "Cloud architecture")

## Messages
All messages are subscribed and published to MQTT broker on topic:
```bash
nes/<username>/message
```

So each client subscribe its own topic on the MQTT broker with its own username.\
And each client send messages to remote username using the similar topic.

Message published to MQTT broker has such format:

```json
{
    "cipherkey": "<base64-encoded encrypted key>",
    "ciphertext": "<base64-encoded encrypted message>",
    "signature": "<base64-encoded message signature>"
}
```

Encrypted message has followed format:
```json
{
    "from": "<recipient username>",
    "to": "<sender username>",
    "timestamp": 112233445566,
    "message": "<text message>"
}
```

## Usage
To use NES Messenger You need to register Your public key to some public keys provider service.
At this moment NES support only PubKey Service, which was created specifically for the needs of NES Messenger.
The PubKey Service is also an open source project, that you can freely deploy to any server.
It doesn't store any private data, except email which is used to register confirmation only.

Source code and more information about PubKey Service are available here: https://github.com/marcinbor85/pubkey.
At this moment it is running on https://microshell.pl/pubkey domain as a default provider for NES.

NES Messenger is a CLI tool, running on Linux and Windows natively thanks to GO. Functionalities are divided into commands,
and the help for each command is available independently with the ```-h``` flag.

```bash
usage: nes <Command> [-h|--help] [-b|--broker "<value>"] [-p|--provider
           "<value>"] [-k|--private "<value>"] [-K|--public "<value>"]
           [-u|--user "<value>"] [-c|--config "<value>"]

           NES messenger

Commands:

  register  Register username at PubKey Service
  listen    Listen to messages
  send      Send message to recipient
  config    Configuration management
  generate  Generate private and public keys pair
  chat      Interactive chat with recipient
  version   Application version

Arguments:

  -h  --help      Print help information
  -b  --broker    MQTT broker server address. Default:
                  tcp://test.mosquitto.org:1883
  -p  --provider  Public key provider address. Default:
                  https://microshell.pl/pubkey
  -k  --private   Private key file. Default: ~/.nes/<user>-rsa
  -K  --public    Public key file. Default: ~/.nes/<user>-rsa.pub
  -u  --user      Local username. Default: <os_user>
  -c  --config    Optional config file. Supported fields: MQTT_BROKER_ADDRESS,
                  PUBKEY_ADDRESS, PRIVATE_KEY_FILE, PUBLIC_KEY_FILE, USERNAME.
                  Default: ~/.nes/config
```

## Examples

### Generate keys
- generate RSA keys pair for current user, and save it to ~/.nes/ directory
```bash
nes generate
```

- generate RSA keys pair with 4096-bits key size for <username>, and save it to ~/.nes/ directory
```bash
nes generate -s 4096 -u <username>
```

- generate RSA keys pair for current user, and save it to provided files
```bash
nes generate -k <private_key_filename> -K <public_key_filename>
```

### Register username
- register current user at public keys provider service
```bash
nes register -e <email>
```

- register provided username with specified public key
```bash
nes register -u <username> -e <email> -K <public_key_filename>
```

### Message listening
- start listening for all messages
```bash
nes listen
```

- start listening for all messages as a specified username
```bash
nes listen -u <username>
```

### Sending messages
- send message to other user
```bash
nes send -t <recipient_username> -m <message>
```

- send message as a specified username to other user
```bash
nes send -t <recipient_username> -u <local_username> -m <message>
```

### Interactive chat
- open interactive chat with user
```bash
nes chat -t <recipient_username>
```

![alt text](assets/chat.png?raw=true "Interactive chat")

- open interactive chat with user using specified configuration file
```bash
nes chat -t <recipient_username> -c <config_file>
```

### Configuration
- save non-volatile settings
```bash
nes config -u <local_username> -S
```

- show current settings
```bash
nes config -s
```

## Contribution
There are a lot of ideas to implement and expand this project.
If you feel that You could participating, don't hesitate for a moment.
It's great fun and you can learn a lot.

If you want to support the project with your work, use the pull request feature.\
If you want to donate (for example for the maintenance of the server providing public key exchange), you can do it here:\
[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donate_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=ZEAEAXGRVZJR8)