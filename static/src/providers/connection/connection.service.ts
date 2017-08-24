import { Injectable } from '@angular/core';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import { CommonService } from '../common/common.service';

@Injectable()
export class ConnectionService {
  private base_url = 'http://127.0.0.1:7410/api/connections/';

  constructor(private common: CommonService) {
  }

  getAllConnections() {
    return this.common.handleGet(this.base_url + 'get_all');
  }

  addConnection(data: FormData) {
    return this.common.handlePost(this.base_url + 'add', data);
  }

  removeConnection(data: FormData) {
    return this.common.handlePost(this.base_url + 'remove', data);
  }
}
