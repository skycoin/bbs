import { Injectable } from '@angular/core';
import { CommonService } from '../common/common.service';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

@Injectable()
export class ApiService {
  private baseUrl = 'http://127.0.0.1:7410/api/';
  private subscriptionsUrl = this.baseUrl + 'subscriptions/';
  private boardsUrl = this.baseUrl + 'content/';
  private submissionAddressUrl = this.boardsUrl + 'meta/submission_addresses/';
  private threadsUrl = this.baseUrl + 'threads/';
  private postsUrl = this.baseUrl + 'posts/';

  constructor(private common: CommonService) {
  }

  addUserVote(data: FormData) {
    return this.common.handlePost(this.baseUrl + 'users/votes/add', data);
  }

  addThreadVote(data: FormData) {
    return this.common.handlePost(this.baseUrl + 'threads/votes/add', data);
  }

  addPostVote(data: FormData) {
    return this.common.handlePost(this.baseUrl + 'posts/votes/add', data);
  }

  getUserVotes(data: FormData) {
    return this.common.handlePost(this.baseUrl + 'threads/votes/get', data);
  }

  getThreadVotes(data: FormData) {
    return this.common.handlePost(this.baseUrl + 'threads/votes/get', data);
  }
  getPostVotes(data: FormData) {
    return this.common.handlePost(this.baseUrl + 'posts/votes/get', data);
  }
  newSeed() {
    return this.common.handleGet(this.baseUrl + 'tools/new_seed');
  }

  getSubmissionAddresses(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'get_all', data);
  }

  addSubmissionAddress(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'add', data);
  }

  removeSubmissionAddress(data: FormData) {
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
    return this.common.handleGet(this.baseUrl + 'node/stats');
  }

  getThreads(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'get_all', data);
  }

  getBoards() {
    return this.common.handleGet(this.boardsUrl + 'get_boards');
  }

  getPosts(data: FormData) {
    return this.common.handlePost(this.postsUrl + 'get_all', data);

  }

  getBoardPage(data: any) {
    return this.common.handlePost(this.boardsUrl + 'get_board_page', data);
  }

  getThreadpage(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'page/get', data);

  }

  addBoard(data: FormData) {
    return this.common.handlePost(this.boardsUrl + 'new_board', data);
  }

  delBoard(data: FormData) {
    return this.common.handlePost(this.boardsUrl + 'delete_board', data);
  }
  addThread(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'new', data);

  }

  addPost(data: FormData) {
    return this.common.handlePost(this.postsUrl + 'add', data);

  }

  importThread(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'import', data);

  }
}
