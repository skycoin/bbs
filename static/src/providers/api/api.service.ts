import { Injectable } from '@angular/core';
import { Http } from "@angular/http";
import { Board, Thread, ThreadPage, Post } from "./msg";

@Injectable()
export class ApiService {
  private base_url = 'http://127.0.0.1:7410/api/'
  constructor(private http: Http) {

  }
  getThreads(key: string) {
    let data = new FormData();
    data.append('board', key);
    return this.handlePost(this.base_url + 'get_threads', data);
  }

  getBoards() {
    return this.handleGet(this.base_url + 'get_boards');
  }

  getPosts(masterKey: string, sub: string) {
    let form = new FormData();
    form.append('board', masterKey);
    form.append('thread', sub);
    return this.handlePost(this.base_url + 'get_posts', form);
  }


  // get_threadpage
  getThreadpage(masterKey: string, sub: string) {
    let form = new FormData();
    form.append('board', masterKey);
    form.append('thread', sub);
    return this.handlePost(this.base_url + 'get_threadpage', form);
  }

  addBoard(data: FormData) {
    return this.handlePost(this.base_url + 'new_board', data);
  }

  addThread(data: FormData) {
    return this.handlePost(this.base_url + 'new_thread', data);
  }

  addPost(data: FormData) {
    return this.handlePost(this.base_url + 'new_post', data);
  }

  changeThread(data: FormData) {
    return this.handlePost(this.base_url + 'import_thread', data);
  }

  private handlePost(url, data: FormData) {
    return new Promise((resolve, reject) => {
      this.http.post(url, data).subscribe(res => {
        resolve(res.json());
      }, err => {
        reject(err);
      })
    });
  }

  private handleGet(url) {
    return new Promise((resolve, reject) => {
      this.http.get(url).subscribe(res => {
        resolve(res.json());
      }, err => {
        reject(err);
      })
    });
  }
}
