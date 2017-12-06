#!/usr/bin/env bash

[[ -f package-lock.json ]] && rm -rf package-lock.json
npm install
npm run build
