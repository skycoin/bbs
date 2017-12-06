import { Injectable } from '@angular/core';
import { CommonService } from '../common/common.service';
import { UserService } from '../user/user.service';
import { LoginSessionUser, PrepareRes } from './msg';
import { Observable } from 'rxjs/Observable';
import { environment } from './../../environments/environment';
import 'rxjs/add/operator/mergeMap';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/of'

@Injectable()
export class ApiService {
  static userInfo: LoginSessionUser = null;
  static sig = '';
  version = 5;
  threadType = 'thread';
  postType = 'post';
  threadVoteType = 'thread_vote';
  postVoteType = 'post_vote';
  userVoteType = 'user_vote';
  private baseUrl = '/api/';
  private adminUrl = this.baseUrl + 'admin/'
  private contentUrl = this.adminUrl + 'content/';
  private sessionUrl = this.adminUrl + 'session/';
  private submissionAddressUrl = this.adminUrl + 'board/';
  private connUrl = this.adminUrl + 'connections/';
  private toolUrl = this.baseUrl + 'tools/'
  private submissionUrl = this.baseUrl + 'new_submission';

  private voteUrl = this.baseUrl + 'submission/';
  private userUrl = this.sessionUrl + 'users/';
  private subscriptionsUrl = this.adminUrl + 'subscriptions/';
  private threadsUrl = this.baseUrl + 'threads/';
  private postsUrl = this.baseUrl + 'posts/';
  constructor(private common: CommonService, private user: UserService) {
    if (environment.production) {
      this.baseUrl = '127.0.0.1:7410/api/';
    }
  }

  // Board
  addBoard(data: FormData) {
    return this.common.handlePost(this.contentUrl + 'new_board', data);
  }

  delBoard(data: FormData) {
    return this.common.handlePost(this.contentUrl + 'delete_board', data);
  }

  // Tools
  newSeed() {
    return this.common.handleGet(this.toolUrl + 'new_seed');
  }

  newKeyPair(data: FormData) {
    return this.common.handlePost(this.toolUrl + 'new_key_pair', data);
  }

  hash(jsonStr: string) {
    const data = new FormData;
    data.append('data', jsonStr);
    return this.common.handlePost(this.toolUrl + 'hash_string', data);
  }

  sig(data: FormData) {
    return this.common.handlePost(this.toolUrl + 'sign', data);
  }

  hashAndSign(jsonStr, secret_key: string) {
    return this.hash(jsonStr).mergeMap(hashData => {
      if (hashData.okay) {
        const data = new FormData;
        data.append('hash', hashData.data.hash);
        data.append('secret_key', secret_key);
        return this.sig(data);
      }
    })
  }

  // Subscriptions
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
    return this.common.handlePost(this.subscriptionsUrl + 'delete', data);
  }

  // Thread
  getBoardPage(data: any) {
    return this.common.handlePost(this.baseUrl + 'get_board_page', data);
  }

  submit(jsonStr, secret_key) {
    return this.hashAndSign(jsonStr, secret_key).mergeMap(signData => {
      console.log('signData:', signData);
      if (signData.okay) {
        const data = new FormData();
        // data.append('type', this.version + ',' + action);
        data.append('body', jsonStr);
        data.append('sig', signData.data.sig);
        return this.common.handlePost(this.submissionUrl, data);
      } else {
        return Observable.throw('Signature Error，please login and try again')
      }
    })
  }

  prepare(action: string, data: FormData) {
    return this.common.handlePost(`${this.baseUrl}submission/prepare_${action}`, data)
  }

  finalize(data: FormData) {
    return this.common.handlePost(`${this.baseUrl}submission/finalize`, data)
  }
  submission(action, data: FormData) {
    return this.prepare(action, data).mergeMap((res: PrepareRes) => {
      const hash = res.data.hash;
      return this.user.hash(res.data.raw).mergeMap(tmpHash => {
        return this.user.sig(tmpHash, this.user.loginInfo.SecKey).mergeMap(sig => {
          data = new FormData()
          data.append('hash', hash)
          data.append('sig', sig)
          return this.finalize(data);
        })
      })
    })
  }
  newThread(jsonStr: string, secret_key: string) {
    return this.submit(jsonStr, secret_key);
  }
  // Thread Page
  getThreadpage(data: FormData) {
    return this.common.handlePost(this.baseUrl + 'get_thread_page', data);
  }
  newPost(jsonStr, secret_key: string) {
    return this.submit(jsonStr, secret_key);
  }
  addThreadVote(jsonStr, secret_key: string) {
    return this.submit(jsonStr, secret_key);
  }
  addPostVote(jsonStr, secret_key: string) {
    return this.submit(jsonStr, secret_key);
  }
  // Other
  newUser(data: FormData) {
    return this.common.handlePost(this.userUrl + 'new', data);
  }
  getFollowPage(data: FormData) {
    return this.common.handlePost(this.baseUrl + 'get_follow_page', data);
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

  addOldThreadVote(data: FormData) {
    return this.common.handlePost(this.voteUrl + 'vote_thread', data);
  }

  addOldPostVote(data: FormData) {
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


  getSubmissionAddresses(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'get_all', data);
  }

  newSubmissionAddress(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'new_submission_address', data);
  }

  delSubmissionAddress(data: FormData) {
    return this.common.handlePost(this.submissionAddressUrl + 'delete_submission_address', data);
  }
  getStats() {
    return this.common.handleGet(this.adminUrl + 'stats');
  }

  getThreads(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'get_all', data);
  }

  getBoards() {
    return this.common.handleGet(this.baseUrl + 'get_boards');
  }

  getPosts(data: FormData) {
    return this.common.handlePost(this.postsUrl + 'get_all', data);

  }
  importThread(data: FormData) {
    return this.common.handlePost(this.threadsUrl + 'import', data);

  }
}
