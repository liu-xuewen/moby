## docker build .
```shell script
DEBU[2020-08-08T21:41:30.468126840+08:00] Calling POST /v1.40/build?buildargs=%7B%7D&cachefrom=%5B%5D&cgroupparent=&cpuperiod=0&cpuquota=0&cpusetcpus=&cpusetmems=&cpushares=0&dockerfile=Dockerfile&labels=%7B%7D&memory=0&memswap=0&networkmode=default&rm=1&shmsize=0&target=&ulimits=null&version=1 
DEBU[2020-08-08T21:41:30.544863281+08:00] Trying to pull alpine from https://registry-1.docker.io v2 
DEBU[2020-08-08T21:41:33.702496832+08:00] Pulling ref from V2 registry: alpine:3.8     
DEBU[2020-08-08T21:41:33.702563662+08:00] docker.io/library/alpine:3.8 resolved to a manifestList object with 6 entries; looking for a unknown/amd64 match 
DEBU[2020-08-08T21:41:33.702581950+08:00] found match for linux/amd64 with media type application/vnd.docker.distribution.manifest.v2+json, digest sha256:954b378c375d852eb3c63ab88978f640b4348b01c1b3456a024a81536dafbbf4 
DEBU[2020-08-08T21:41:35.804659563+08:00] pulling blob "sha256:486039affc0ad0f17f473efe8fb25c947515a8929198879d1e64210ef142372f" 
DEBU[2020-08-08T21:41:37.076614357+08:00] Downloaded 486039affc0a to tempfile /var/lib/docker/tmp/GetImageBlob080586704 
DEBU[2020-08-08T21:41:37.078089924+08:00] Applying tar in /var/lib/docker/overlay2/aceddbe1b592d480f76d91378657867a328476aa22902062844a6c03d1fae248/diff  storage-driver=overlay2
DEBU[2020-08-08T21:41:37.266393276+08:00] Applied tar sha256:7444ea29e45e927abea1f923bf24cac20deaddea603c4bb1c7f2f5819773d453 to aceddbe1b592d480f76d91378657867a328476aa22902062844a6c03d1fae248, size: 4413305 
DEBU[2020-08-08T21:41:37.311157918+08:00] [BUILDER] Cache miss: [/bin/sh -c #(nop)  MAINTAINER liiuxuewen@gmail.com] 
DEBU[2020-08-08T21:41:37.311187243+08:00] [BUILDER] Command to be executed: [/bin/sh -c #(nop)  MAINTAINER liiuxuewen@gmail.com] 
DEBU[2020-08-08T21:41:37.330683379+08:00] container mounted via layerStore: &{/var/lib/docker/overlay2/c4f82f9b1998ef0933dc90c7ea278c2b5821489c6f62cf0365b26f1d91d8e1b8/merged 0x7f56957f3ac0 0x7f56957f3ac0} 
DEBU[2020-08-08T21:41:37.347974270+08:00] Applying tar in /var/lib/docker/overlay2/bfce382dbc5d82f4b8286fe1e64c4cc756a73eacede7f238efab9b4cdbab18cc/diff  storage-driver=overlay2
DEBU[2020-08-08T21:41:37.408416119+08:00] Applied tar sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef to bfce382dbc5d82f4b8286fe1e64c4cc756a73eacede7f238efab9b4cdbab18cc, size: 0 
INFO[2020-08-08T21:41:37.440492917+08:00] Layer sha256:d976c902143a892da7547e763e06c4dea0bea71f32d55ea542edcc9595f7a2ac cleaned up 
DEBU[2020-08-08T21:41:38.009793354+08:00] CopyFileWithTar(/var/lib/docker/tmp/docker-builder726125813/unix-sock, /var/lib/docker/overlay2/8a7bbae79767167a97e4f4026605c03fcb0a0be895df2033b395b5511b976a6a/merged/code/unix-sock) 
DEBU[2020-08-08T21:41:38.106377698+08:00] Applying tar in /var/lib/docker/overlay2/e21b8cc335663795b14a27e60c697ff6aa9e10a451c8961c83a209109209f04c/diff  storage-driver=overlay2
DEBU[2020-08-08T21:41:38.200452039+08:00] Applied tar sha256:9bf7999d4fe71c5f79ffa88de61ba831b47d80b2671cee56959b6b37b38d049f to e21b8cc335663795b14a27e60c697ff6aa9e10a451c8961c83a209109209f04c, size: 3010780 
DEBU[2020-08-08T21:41:39.011770554+08:00] [BUILDER] Command to be executed: [/bin/sh -c #(nop)  STOPSIGNAL SIGTERM] 
DEBU[2020-08-08T21:41:39.031449991+08:00] container mounted via layerStore: &{/var/lib/docker/overlay2/8d2b41f65d1f7a7c471e631a0482a65df2e4bc14ad083565072486fda0c7aa52/merged 0x7f56957f3ac0 0x7f56957f3ac0} 
DEBU[2020-08-08T21:41:39.047510461+08:00] Applying tar in /var/lib/docker/overlay2/7327cfa1fb0d1d6570bd4906fcddb7d8b587bcf9e34f23fe41ff3beaf10ca97f/diff  storage-driver=overlay2
DEBU[2020-08-08T21:41:39.120433122+08:00] Applied tar sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef to 7327cfa1fb0d1d6570bd4906fcddb7d8b587bcf9e34f23fe41ff3beaf10ca97f, size: 0 
INFO[2020-08-08T21:41:39.147089379+08:00] Layer sha256:8e10700d131068cf0a025d398dc20b95f99e907a606a967a0703765e412c991f cleaned up 

```

#### Dockerfile
```shell script
FROM alpine:3.8

MAINTAINER liiuxuewen@gmail.com

COPY --chown=nobody:nobody ./unix-sock /code/unix-sock

STOPSIGNAL SIGTERM

```
