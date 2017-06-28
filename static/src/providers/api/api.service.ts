import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Board, Thread, ThreadPage, Post } from './msg';
import { CommonService } from '../common/common.service';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

@Injectable()
export class ApiService {
  private base_url = 'http://127.0.0.1:7410/api/'
  private submissionAddressUrl = this.base_url + 'boardmeta/';
  constructor(private http: Http, private common: CommonService) {
  }
  getSubmissionAddresses(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'get_submissionaddresses', data);
  }
  addSubmissionAddress(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'add_submissionaddress', data);
  }
  removeSubmissionAddress(data:FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'remove_submissionaddress', data);
  }
  getSubscriptions() {
    return this.common.handleGet(this.base_url + 'get_subscriptions');
  }

  getSubscription(data: FormData) {
    return this.common.handlePost(this.base_url + 'get_subscription', data);
  }
  subscribe(data: FormData) {
    return this.common.handlePost(this.base_url + 'subscribe', data);
  }
  unSubscribe(data: FormData) {
    return this.common.handlePost(this.base_url + 'unsubscribe', data);
  }
  getStats() {
    return this.common.handleGet(this.base_url + 'get_stats');
  }

  getThreads(data: FormData) {
    return this.common.handlePost(this.base_url + 'get_threads', data);
  }

  getBoards() {
    return this.common.handleGet(this.base_url + 'get_boards');
  }

  getPosts(data: FormData) {
    return this.common.handlePost(this.base_url + 'get_posts', data);

  }

  getBoardPage(data: FormData) {
    return this.common.handlePost(this.base_url + 'get_boardpage', data);
  }
  getThreadpage(data: FormData) {
    return this.common.handlePost(this.base_url + 'get_threadpage', data);

  }

  addBoard(data: FormData) {
    return this.common.handlePost(this.base_url + 'new_board', data);

  }

  addThread(data: FormData) {
    return this.common.handlePost(this.base_url + 'new_thread', data);

  }

  addPost(data: FormData) {
    return this.common.handlePost(this.base_url + 'new_post', data);

  }

  importThread(data: FormData) {
    return this.common.handlePost(this.base_url + 'import_thread', data);

  }
}
