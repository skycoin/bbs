import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

import { CommonService } from "../common/common.service";
import { User } from "./user.msg";

@Injectable()
export class UserService {
  private base_url = 'http://127.0.0.1:7410/api/users/'
  constructor(private http: Http, private common: CommonService) { }


  getCurrent() {
    return this.http.get(this.base_url + 'get_current')
      .map((response: Response) => <User>response.json())
      .catch(err => this.common.handleError(err));
  }

  getAllMasters() {
    return this.http.get(this.base_url + 'get_masters').map((res: Response) => <Array<User>>res.json()).catch(err => this.common.handleError(err));
  }

  getAll() {
    return this.http.get(this.base_url + 'get_all').map((res: Response) => <Array<User>>res.json()).catch(err => this.common.handleError(err));
  }

  setCurrent(data: FormData) {
    return this.http.post(this.base_url + 'set_current', data).map((res: Response) => <User>res.json()).catch(err => this.common.handleError(err));
  }

  newMaster(data: FormData) {
    return this.http.post(this.base_url + 'new_master', data).map((res: Response) => <User>res.json()).catch(err => this.common.handleError(err));
  }

  //not master 
  newUser(data: FormData) {
    return this.http.post(this.base_url + 'new', data).map((res: Response) => <User>res.json()).catch(err => this.common.handleError(err));
  }

  remove(data: FormData) {
    return this.http.post(this.base_url + 'remove', data).map((res: Response) => res.json()).catch(err => this.common.handleError(err));
  }
}