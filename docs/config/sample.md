# Sample Config { data-search-exclude }

A full config file template has been provided below:

```yaml
app:
  mode:
  api:
    enabled:
    port:
    
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
    snap_bbox:
    snap_timestamp:
    snap_crop:
    snap_hires:
    max_snap_retry:
    notify_once:
    notify_detections:
    recheck_delay:
    audio_only:

  quiet:
    start:
    end:

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

  license_plate:
    enabled: false
    allow:
    block:

  apprise_api:
    enabled: false
    server:
    token:
    urls:
    tags:
    ignoressl:
    template:

  discord:
    enabled: false
    webhook:
    disable_embed:
    template:
  
  gotify:
    enabled: false
    server:
    token: 
    ignoressl:
    template:
  
  matrix:
    enabled: false
    server:
    username:
    password:
    roomid:
    ignoressl:
    template:

  mattermost:
    enabled: false
    webhook:
    channel:
    username:
    priority:
    ignoressl:
    headers:
    template:

  ntfy:
    enabled: false
    server:
    topic:
    ignoressl:
    headers:
    template:

  pushover:
    enabled: false
    token:
    userkey:
    devices:
    sound:
    priority:
    retry:
    expire:
    ttl:
    template:

  signal:
    enabled: false
    server:
    account:
    recipients:
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
    send_clip:
    message_thread_id:
    token:
    template:

  webhook:
    enabled: false
    server: 
    ignoressl:
    method:
    params:
    headers:
    template:


monitor:
  enabled: false
  url: 
  interval: 
  ignoressl: 
```
