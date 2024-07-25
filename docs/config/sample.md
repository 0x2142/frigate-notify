# Sample Config { data-search-exclude }

A full config file template has been provided below:

```yaml
frigate:
  server: 
  ignoressl:
  public_url:
  headers:
  startup_check:
    attempts:
    interval:

  webapi:
    enabled: 
    interval: 
    
  mqtt: 
    enabled: 
    server: 
    port:
    clientid:
    username: 
    password: 
    topic_prefix: 
  
  cameras:
    exclude:
      - test_cam_01

alerts:  
  general:
    title:
    timeformat:
    nosnap:

  zones:
    unzoned: allow
    allow:
     - test_zone_01
    block:
     - test_zone_02

  labels:
    min_score:
    allow:
     - person
     - dog
    block:
     - bird

  sublabels:
    allow:
     - ABCD
    block:
     - EFGH

  discord:
    enabled: false
    webhook: 
    template:
  
  gotify:
    enabled: false
    server:
    token: 
    ignoressl:
    template:
  
  smtp:
    enabled: false
    server: 
    port: 
    tls: 
    user:
    password:
    recipient:
    template:

  telegram:
    enabled: false
    chatid:
    token:
    template:

  pushover:
    enabled: false
    token:
    userkey:
    devices:
    priority:
    retry:
    expire:
    ttl:
    template:

  ntfy:
    enabled: false
    server:
    topic:
    ignoressl:
    headers:
    template:

  webhook:
    enabled: false
    server: 
    ignoressl:
    headers:
    template:

monitor:
  enabled: false
  url: 
  interval: 
  ignoressl: 
```
