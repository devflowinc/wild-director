### wild-director

`wild-director` runs a DNS server and a HTTP/TLS server to create wildcard
subdomain alias proxies for Trieve demos.

#### usage

```
make binary
make deploy
```

#### configuration

Edit `/etc/wild.json` to configure on the machine to change the list of
aliases.

#### TLS

Every 90 days, run the following `certbot` command to renew certs.

```
certbot -d '*.demo.trytrieve.com' -d 'demo.trytrieve.com' --manual --preferred-challenges dns certonly
```

The output should contain something like 

```
Please deploy a DNS TXT record under the name:

_acme-challenge.demo.trytrieve.com.

with the following value:

dBSfPhqIfk6IVkanIvguzC3Y4oxAOV_EzUeb4KBYcPg
```

Copy the code to `/etc/systemd/system/wild-director.service` like this:

```
Environment="WILD_TXT=dBSfPhqIfk6IVkanIvguzC3Y4oxAOV_EzUeb4KBYcPg"
```
