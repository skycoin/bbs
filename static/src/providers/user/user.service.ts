import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

import { CommonService } from '../common/common.service';
import { User } from './user.msg';

@Injectable()
export class UserService {
  private baseUrl = 'http://127.0.0.1:7410/api/users/';
  private mastersUrl = this.baseUrl + 'masters/';
  private currentUrl = this.mastersUrl + 'current/';

  constructor(private http: Http, private common: CommonService) {
  }


  getCurrent() {
    return this.http.get(this.currentUrl + 'get')
      .map((response: Response) => <User>response.json())
      .catch(err => this.common.handleError(err));
  }

  getAllMasters() {
    return this.http.get(this.mastersUrl + 'get_all').map((res: Response) => <Array<User>>res.json())
      .catch(err => this.common.handleError(err));
  }

  getAll() {
    return this.http.get(this.baseUrl + 'get_all').map((res: Response) => <Array<User>>res.json())
      .catch(err => this.common.handleError(err));
  }

  setCurrent(data: FormData) {
    return this.http.post(this.currentUrl + 'set', data).map((res: Response) => <User>res.json())
      .catch(err => this.common.handleError(err));
  }

  newMaster(data: FormData) {
    return this.http.post(this.mastersUrl + 'add', data).map((res: Response) => <User>res.json())
      .catch(err => this.common.handleError(err));
  }

  newOrModifyUser(data: FormData) {
    return this.http.post(this.baseUrl + 'add', data).map((res: Response) => <User>res.json())
      .catch(err => this.common.handleError(err));
  }

  remove(data: FormData) {
    return this.http.post(this.baseUrl + 'remove', data).map((res: Response) => res.json())
      .catch(err => this.common.handleError(err));
  }
}
