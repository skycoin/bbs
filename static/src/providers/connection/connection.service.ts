import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import { CommonService } from '../common/common.service';

@Injectable()
export class ConnectionService {
  private base_url = 'http://127.0.0.1:7410/api/connections/'
  constructor(private http: Http, private common: CommonService) { }
  getAllConnections() {
    return this.http.get(this.base_url + 'get_all').
      map((res: Response) => res.json()).
      catch(err => this.common.handleError(err));
  }
  addConnection(data: FormData) {
    return this.http.post(this.base_url + 'new', data).
      map((res: Response) => res.json()).
      catch(err => this.common.handleError(err));
  }
  removeConnection(data: FormData) {
    return this.http.post(this.base_url + 'remove', data).
      map((res: Response) => res.json()).
      catch(err => this.common.handleError(err));
  }
} 
