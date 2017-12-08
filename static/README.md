# Skycoin BBS Web UI
The user interface of Skycoin BBS is represented via a locally served website. In the future, users will no longer need to run a full BBS node to access/submit content to boards. And instead, nodes will host public websites.

The project contains both the source (src) and target (dist) files of this web interface.

## Prerequisites

The project have dependencies that require Node 6.9.0 or higher, together
with NPM 3 or higher.

#### Updating npm

You may need to run this with `sudo`.

```bash
npm install npm@latest -g
```

#### Updating/Installing ng

```bash
npm install -g @angular/cli
```

## Run
From your command line:
```bash
git clone https://github.com/skycoin/bbs.git
cd bbs/static
npm install
npm run start
```
or
```bash
git clone https://github.com/skycoin/bbs.git
cd bbs/static
./run.sh
```

## Build
From your command line:
```bash
git clone https://github.com/skycoin/bbs.git
cd bbs/static
npm install
npm run build
```
or
```bash
git clone https://github.com/skycoin/bbs.git
cd bbs/static
./build.sh
```
