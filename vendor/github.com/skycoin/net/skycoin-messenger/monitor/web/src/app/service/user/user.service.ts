import { Injectable } from '@angular/core';

@Injectable()
export class UserService {
  SSHCLIENTINFO = '_SKYWIRE_SSHCLIENTINFO';
  SOCKETCLIENTINFO = '_SKYWIRE_SOCKETCLIENTINFO';
  HOMENODELABLE = '_SKYWIRE_HOMENODELABEL';
  constructor() { }
  saveHomeLabel(nodeKey: string, label: string) {
    let homeLabels = this.get(this.HOMENODELABLE);
    if (!homeLabels) {
      homeLabels = {};
    }
    homeLabels[nodeKey] = label;
    localStorage.setItem(this.HOMENODELABLE, JSON.stringify(homeLabels));
  }
  saveClientConnectInfo(data: ConnectServiceInfo, key: string) {
    let info = <Array<ConnectServiceInfo>>this.get(key);
    if (info) {
      const len = info.length;
      if (len >= 5) {
        info[len - 1] = data;
      } else {
        let isExist = false;
        info.forEach(v => {
          if (v.appKey === data.appKey && v.nodeKey === data.nodeKey) {
            v.count += 1;
            isExist = true;
          }
        });
        if (!isExist) {
          info.push(data);
        }
      }
    } else {
      info = [];
      info.push(data);
    }
    this.sort(info);
    localStorage.setItem(key, JSON.stringify(info));
  }
  removeClientConnectInfo(key: string, index: number) {
    const info = <Array<ConnectServiceInfo>>this.get(key);
    info.splice(index, 1);
    this.sort(info);
    localStorage.setItem(key, JSON.stringify(info));
  }
  editClientConnectInfo(item: ConnectServiceInfo, key: string, index: number) {
    const info = <Array<ConnectServiceInfo>>this.get(key);
    info[index] = item;
    localStorage.setItem(key, JSON.stringify(info));
  }
  get(key: string) {
    return JSON.parse(localStorage.getItem(key));
  }
  sort(info: Array<ConnectServiceInfo>) {
    info.sort((a: ConnectServiceInfo, b: ConnectServiceInfo) => {
      if (a.count < b.count) {
        return 1;
      }
      if (a.count > b.count) {
        return -1;
      }
      return a.nodeKey.localeCompare(b.nodeKey);
    });
  }
}
// export interface LinkedList
export interface ConnectServiceInfo {
  label?: string;
  nodeKey?: string;
  appKey?: string;
  count?: number;
}
