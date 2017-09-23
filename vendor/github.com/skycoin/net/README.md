# Skycoin Networking Framework



## Protocol

```
                  +--+--------+--------+--------------------+
msg protocol      |  |        |        |                    |
                  +-++-------++-------++---------+----------+
                    |        |        |          |
                    v        |        v          v
                  msg type   |     msg len    msg body
                   1 byte    v     4 bytes
                          msg seq
                          4 bytes



                  +-----------+--------+--------------------+
normal msg        |01|  seq   |  len   |       body         |
                  +-----------+--------+--------------------+


                  +-----------+
ack msg           |80|  seq   |
                  +-----------+


                  +--------------------+
ping msg          |81|    timestamp    |
                  +--------------------+


                  +--------------------+
pong msg          |82|    timestamp    |
                  +--------------------+
```
