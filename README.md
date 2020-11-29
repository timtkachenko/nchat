# nchat

## build

`make nchat`

`make nclient`

## run

nchat [CHAT_PORT] [NANOMSG_NODE_ADDR]

**default server port 22111**

starts first node on 22111 and nanomsg address 0.0.0.0:22112
```sh
$ ./nchat
```
to join first node run with
```sh
$ ./nchat 22333 0.0.0.0:22112
```

nclient [CHAT_ADDR] [USER_ID] [MANDATORY_FRIEND_ID]

connect clients
```sh
$ ./nclient 0.0.0.0:22111 12 1
```

```sh
$ ./nclient 0.0.0.0:22333 1 12
```
