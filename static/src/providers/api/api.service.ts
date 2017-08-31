import { Injectable } from '@angular/core';
import { CommonService } from '../common/common.service';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

@Injectable()
export class ApiService {
  private baseUrl = 'http://127.0.0.1:7410/api/';
  private voteUrl = this.baseUrl + 'votes/';
  private connUrl = this.baseUrl + 'connections/'
  private sessionUrl = this.baseUrl + 'session/'
  private userUrl = this.sessionUrl + 'users/';
  private subscriptionsUrl = this.baseUrl + 'subscriptions/';
  private contentUrl = this.baseUrl + 'content/';
  private submissionAddressUrl = this.baseUrl + 'admin/board/';
  private threadsUrl = this.baseUrl + 'threads/';
  private postsUrl = this.baseUrl + 'posts/';

  constructor(private common: CommonService) {
  }
  newUser(data: FormData) {
    return this.common.handlePost(this.userUrl + 'new', data);
  }
  getFollowPage(data: FormData) {
    return this.common.handlePost(this.voteUrl + 'get_follow_page', data);
  }
  delConnection(data: FormData) {
    return this.common.handlePost(this.connUrl + 'delete', data);
  }
  newConnection(data: FormData) {
    return this.common.handlePost(this.connUrl + 'new', data);
  }
  getAllConnections() {
    return this.common.handleGet(this.connUrl + 'get_all');
  }
  delUser(data: FormData) {
    return this.common.handlePost(this.userUrl + 'delete', data);
  }
  getAllUser() {
    return this.common.handleGet(this.userUrl + 'get_all');

  }
  logout() {
    return this.common.handleGet(this.sessionUrl + 'logout');
  }
  login(data: FormData) {
    return this.common.handlePost(this.sessionUrl + 'login', data);
  }
  getSessionInfo() {
    return this.common.handleGet(this.sessionUrl + 'get_info');
  }

  addUserVote(data: FormData) {
    return this.common.handlePost(this.voteUrl + 'vote_user', data);
  }

  addThreadVote(data: FormData) {
    return this.common.handlePost(this.voteUrl + 'vote_thread', data);
  }

  addPostVote(data: FormData) {
    return this.common.handlePost(this.voteUrl + 'vote_post', data);
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

  newSubmissionAddress(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'new_submission_address', data);
  }

  delSubmissionAddress(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'delete_submission_address', data);
  }

  getSubscriptions() {
    return this.common.handleGet(this.subscriptionsUrl + 'get_all');
  }

  getSubscription(data: FormData) {
    return this.common.handlePost(this.subscriptionsUrl + 'get', data);
  }

  newSubscription(data: FormData) {
    return this.common.handlePost(this.subscriptionsUrl + 'new', data);
  }

  delSubscription(data: FormData) {
    return this.common.handlePost(this.subscriptionsUrl + 'remove', data);
  }

  getStats() {
    return this.common.handleGet(this.baseUrl + 'node/stats');
  }

  getThreads(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'get_all', data);
  }

  getBoards() {
    return this.common.handleGet(this.contentUrl + 'get_boards');
  }

  getPosts(data: FormData) {
    return this.common.handlePost(this.postsUrl + 'get_all', data);

  }

  getBoardPage(data: any) {
    return this.common.handlePost(this.contentUrl + 'get_board_page', data);
  }

  getThreadpage(data: FormData) {
    return this.common.handlePost(this.contentUrl + 'get_thread_page', data);

  }

  addBoard(data: FormData) {
    return this.common.handlePost(this.contentUrl + 'new_board', data);
  }

  delBoard(data: FormData) {
    return this.common.handlePost(this.contentUrl + 'delete_board', data);
  }
  newThread(data: FormData) {
    return this.common.handlePost(this.contentUrl + 'new_thread', data);

  }

  newPost(data: FormData) {
    return this.common.handlePost(this.contentUrl + 'new_post', data);

  }

  importThread(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'import', data);

  }
}
