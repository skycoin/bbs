import {Injectable} from "@angular/core";
import {Http} from "@angular/http";
import {CommonService} from "../common/common.service";
import "rxjs/add/operator/map";
import "rxjs/add/operator/catch";

@Injectable()
export class ApiService {
  private baseUrl = 'http://127.0.0.1:7410/api/';
  private statsUrl = this.baseUrl + 'stats/';
  private subscriptionsUrl = this.baseUrl + 'subscriptions/';
  private boardsUrl = this.baseUrl + 'boards/';
  private submissionAddressUrl = this.boardsUrl + 'meta/submission_addresses/';
  private threadsUrl = this.baseUrl + 'threads/';
  private postsUrl = this.baseUrl + 'posts/';
  constructor(private http: Http, private common: CommonService) {
  }
  getSubmissionAddresses(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'get_all', data);
  }
  addSubmissionAddress(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'add', data);
  }
  removeSubmissionAddress(data:FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'remove', data);
  }
  getSubscriptions() {
    return this.common.handleGet(this.subscriptionsUrl + 'get_all');
  }

  getSubscription(data: FormData) {
    return this.common.handlePost(this.subscriptionsUrl + 'get', data);
  }
  subscribe(data: FormData) {
    return this.common.handlePost(this.subscriptionsUrl + 'add', data);
  }
  unSubscribe(data: FormData) {
    return this.common.handlePost(this.subscriptionsUrl + 'remove', data);
  }
  getStats() {
    return this.common.handleGet(this.statsUrl + 'get');
  }

  getThreads(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'get_all', data);
  }

  getBoards() {
    return this.common.handleGet(this.boardsUrl + 'get_all');
  }

  getPosts(data: FormData) {
    return this.common.handlePost(this.postsUrl + 'get_all', data);

  }

  getBoardPage(data: FormData) {
    return this.common.handlePost(this.boardsUrl + 'page/get', data);
  }
  getThreadpage(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'page/get', data);

  }

  addBoard(data: FormData) {
    return this.common.handlePost(this.boardsUrl + 'add', data);

  }

  addThread(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'add', data);

  }

  addPost(data: FormData) {
    return this.common.handlePost(this.postsUrl + 'add', data);

  }

  importThread(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'import', data);

  }
}
