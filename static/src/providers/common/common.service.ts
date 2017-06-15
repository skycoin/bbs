import { Http, Response } from "@angular/http";
import { Injectable } from '@angular/core';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/filter';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class CommonService {
  private alertType: string = 'info';
  private alertMessage: string = '';
  private alert: boolean = false;
  topBtn: boolean = false;
  loading: boolean = false;
  constructor(private http: Http) { }

  handleError(error: Response) {
    console.error('Error:', error.json() || 'Server error', 'danger');
    this.showAlert((error.json() instanceof Object ? 'Server error': error.json()) || 'Server error', 'danger', 3000);
    return Observable.throw(error.json() || 'Server error');
  }
  handleGet(url: string) {
    if (!url) {
      return Observable.throw('The connection is empty');
    }
    return this.http.get(url).filter((res: Response) => res.status === 200).map((res: Response) => res.json()).catch(err => this.handleError(err));
  }

  handlePost(url: string, data: FormData) {
    if (!url || !data) {
      return Observable.throw('Parameters and connections can not be empty');
    }
    return this.http.post(url, data).filter((res: Response) => res.status === 200).map((res: Response) => res.json()).catch(err => this.handleError(err));
  }

  showAlert(message: string, type?: string, timeout?: number) {
    this.alert = false;
    this.alertMessage = message;
    if (type) {
      this.alertType = type;
    }
    if (timeout > 0) {
      setTimeout(() => {
        this.alert = false;
      }, timeout);
    }
    this.alert = true;
  }

  showOrHideToTopBtn() {
    let pos = (document.documentElement.scrollTop || document.body.scrollTop) + document.documentElement.offsetHeight;
    let max = document.documentElement.scrollHeight;
    if (pos > (max / 3)) {
      this.topBtn = true;
    } else {
      this.topBtn = false;
    }
  }

  scrollToTop() {
    window.scrollTo(0, 0);
  }
}