import { Injectable } from '@angular/core';
import { Http, Response } from "@angular/http";
import { Board, Thread, ThreadPage, Post } from "./msg";
import { CommonService } from "../common/common.service";
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

@Injectable()
export class ApiService {
  private base_url = 'http://127.0.0.1:7410/api/'
  constructor(private http: Http, private common: CommonService) {

  }
  getThreads(key: string) {
    let data = new FormData();
    data.append('board', key);
    return this.http.post(this.base_url + 'get_threads', data).map((res: Response) => res.json()).catch(err => this.common.handleError(err));
  }

  getBoards() {
    return this.http.get(this.base_url + 'get_boards').map((res: Response) => res.json()).catch(err => this.common.handleError(err));
  }

  getPosts(masterKey: string, sub: string) {
    let data = new FormData();
    data.append('board', masterKey);
    data.append('thread', sub);
    return this.http.post(this.base_url + 'get_posts', data).map((res: Response) => res.json()).catch(err => this.common.handleError(err));

  }


  // get_threadpage
  getThreadpage(masterKey: string, sub: string) {
    let data = new FormData();
    data.append('board', masterKey);
    data.append('thread', sub);
    return this.http.post(this.base_url + 'get_threadpage', data).map((res: Response) => res.json()).catch(err => this.common.handleError(err));

  }

  addBoard(data: FormData) {
    return this.http.post(this.base_url + 'new_board', data).map((res: Response) => res.json()).catch(err => this.common.handleError(err));

  }

  addThread(data: FormData) {
    return this.http.post(this.base_url + 'new_thread', data).map((res: Response) => res.json()).catch(err => this.common.handleError(err));

  }

  addPost(data: FormData) {
    return this.http.post(this.base_url + 'new_post', data).map((res: Response) => res.json()).catch(err => this.common.handleError(err));

  }

  importThread(data: FormData) {
    return this.http.post(this.base_url + 'import_thread', data).map((res: Response) => res.json()).catch(err => this.common.handleError(err));

  }
}
