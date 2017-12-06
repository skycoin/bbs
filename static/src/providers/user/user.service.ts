import { Injectable } from '@angular/core';
import * as crypto from 'crypto-js';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/of'
declare var cipher: any;

@Injectable()
export class UserService {
  // private _bbsKey = '';
  private rootKey = 'skybbs_users';
  private tmpUser = 'tmp_user';
  private sessionUser = 'session_user';
  loginInfo: { PublicKey?: string, SecKey?: string, Seed?: string } = null;
  constructor() {}

  setTmpItem(name: string) {
    const item = { name: name, timestamp: new Date().getTime() + (86400000 * 7) }
    localStorage.setItem(this.tmpUser, JSON.stringify(item));
    // sessionStorage.setItem(this.sessionUser, JSON.stringify(data));
  }

  getTmpItem() {
    const item = JSON.parse(localStorage.getItem(this.tmpUser));
    if (item) {
      if (item['timestamp'] < new Date().getTime()) {
        localStorage.removeItem(this.tmpUser);
        return null;
      }
      const info = JSON.parse(sessionStorage.getItem(this.sessionUser))
      if (info) {
        item.data = info;
      }
      return item;
    }
    return null;
  }

  setItem(key: string, data: any) {
    let orginData = this.getOrginData();
    if (!orginData) {
      orginData = {};
    }
    orginData[key] = data;
    localStorage.setItem(this.rootKey, JSON.stringify(orginData));
  }
  getItem(key) {
    const data = this.getOrginData();
    return data[key];
  }
  getUserList() {
    const data = this.getOrginData();
    const list = [];
    // tslint:disable-next-line:forin
    for (const key in data) {
      list.push(key);
    }
    return list;
  }
  getOrginData(): any {
    return JSON.parse(localStorage.getItem(this.rootKey));
  }
  newSeed() {
    return cipher.generateSeed();
  }
  newKeyPair(seed: string) {
    return cipher.generateKeyPair(seed);
  }
  hash(data: string) {
    return Observable.of(cipher.hash(data));
  }
  sig(hash: string, secKey: string) {
    return Observable.of(cipher.sig(hash, secKey));
  }
  encrypt(data, password: string) {
    return Observable.of(crypto.AES.encrypt(data, password).toString())
  }
  decrypt(data, password: string) {
    const bytes = crypto.AES.decrypt(data, password);
    if (bytes.words[0] < 0) {
      return Observable.of('')
    }
    const plaintext = bytes.toString(crypto.enc.Utf8);
    return Observable.of(JSON.parse(plaintext))
  }
}

export interface Cipher {
  generateSeed: Function;
  generateKeyPair: Function;
  hash: Function;
  sig: Function;
}
