# pfui


This a simple UI for `pf` to block devices on the LAN.
This works by adding/removing a host from a PF table (`nointernet` in the pf.conf example below).

### configure PF

```pf.conf

# any ip in this table will not be allowed to reach outside
table <nointernet> persist

block in quick from <nointernet>
block out quick to <nointernet>

```

### configure the hosts you want to block in `config.json`:

```json
{
  "table": "nointernet",
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

Run the tool and open the browser at the routers's internal IP:

http://192.168.1.1:3000/psui/devices


### TODO:
    [] - rc file
    [] - embed the html file
    [] - option to ignode the devices from the config file and display all active LAN devices
