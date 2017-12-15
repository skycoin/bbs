cd /tmp

# TODO: wget all stuff.

sudo useradd skycoin -s /sbin/nologin -M

# TODO: cp skycoin-discoverynode.service
sudo chmod 755 /lib/systemd/system/skycoin-discoverynode.service

sudo systemctl enable skycoin-discoverynode.service
sudo systemctl start skycoin-discoverynode
sudo journalctl -f -u skycoin-discoverynode

# TODO: cp 30-skycoin-discoverynode.conf to /etc/rsyslog.d
sudo systemctl restart rsyslog
netstat -an | grep "LISTEN "
sudo systemctl restart skycoin-discoverynode