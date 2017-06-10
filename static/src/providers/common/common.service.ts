import { Http } from "@angular/http";
import { Injectable } from '@angular/core';
import 'rxjs/add/observable/throw';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class CommonService {
  constructor(private http: Http) { }
 
  handleError(error: Response) {
    return Observable.throw(error.json() || 'Server error');
  }
}