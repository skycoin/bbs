# Skycoin BBS Web UI
The user interface of Skycoin BBS is represented via a locally served website. In the future, users will no longer need to run a full BBS node to access/submit content to boards. And instead, nodes will host public websites.

## Prerequisites

Both the CLI and generated project have dependencies that require Node 6.9.0 or higher, together
with NPM 3 or higher.

## Installation


### node and git
You'll need [Git](https://git-scm.com) and [Node.js](https://nodejs.org/en/download/) (which comes with [npm](http://npmjs.com)) installed on your computer.

### ng

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

## Build
From your command line:
```bash
git clone https://github.com/skycoin/bbs.git
cd bbs/static
npm install
npm run build
```
