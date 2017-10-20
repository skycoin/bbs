import { Injectable } from '@angular/core';
import * as crypto from 'crypto-js';

declare var cipher: Cipher;

@Injectable()
export class UserService {
  private rootKey = 'skybbs_users'
  constructor() { }

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
  encrypt(data, password: string) {
    return crypto.AES.encrypt(data, password).toString();
  }
  decrypt(data, password: string) {
    const bytes = crypto.AES.decrypt(data, password);
    const plaintext = bytes.toString(crypto.enc.Utf8);
    return JSON.parse(plaintext);
  }
}

export interface Cipher {
  generateSeed: Function;
  generateKeyPair: Function;
}
