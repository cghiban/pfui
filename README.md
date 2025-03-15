# pfui


This a simple UI for `pf` to block devices on the LAN.
This works by adding/removing a host from a PF table (`nointernet` in the pf.conf example below).

#### configure PF

```pf.conf

# any ip in this table will not be allowed to reach outside
table <nointernet> persist

block in quick from <nointernet>
block out quick to <nointernet>

```

Configure the hosts you want to block in `config.json`:

```json
{
  "table": "nointernet",
  "auth": {
    "user": "xxx",
    "pass": "zzz"
  },
  "devices": [
    {
      "mac": "cc:ee:aa:01:5d:d5",
      "name": "PC"
    },
    {
      "mac": "10:56:aa:83:1f:15",
      "name": "tablet"
    }
  ]
}
```
