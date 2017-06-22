import { Http, Response } from '@angular/http';
import { Injectable } from '@angular/core';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/do';
import { Observable } from 'rxjs/Observable';
import { LoadingComponent } from '../../components';

@Injectable()
export class CommonService {
  private alertType = 'info';
  private alertMessage = '';
  private alert = false;
  topBtn = false;
  loading: LoadingComponent = null;
  constructor(private http: Http) { }

  handleError(error: Response) {
    console.error('Error:', error.json() || 'Server error', 'danger');
    this.showAlert((error.json() instanceof Object ? 'Server error' : error.json()) || 'Server error', 'danger', 3000);
    return Observable.throw(error.json() || 'Server error');
  }
  handleGet(url: string) {
    if (!url) {
      return Observable.throw('The connection is empty');
    }
    // if (this.loading) {
    //   this.loading.start();
    // }
    return this.http.get(url).
      filter((res: Response) => res.status === 200).
      map((res: Response) => res.json()).
      // do(() => { if (this.loading) { this.loading.close() } }).
      catch(err => this.handleError(err));
  }

  handlePost(url: string, data: FormData) {
    if (!url || !data) {
      return Observable.throw('Parameters and connections can not be empty');
    }
    // if (this.loading) {
    //   this.loading.start();
    // }
    return this.http.post(url, data).
      filter((res: Response) => res.status === 200).
      map((res: Response) => res.json()).
      // do(() => { if (this.loading) { this.loading.close() } }).
      catch(err => this.handleError(err));
  }



  /**
   * Show Error Alert
   * @param message Error Text
   * @param timeout
   */
  showErrorAlert(message: string, timeout: number = 3000) {
    this.showAlert(message, 'danger', timeout);
  }
  showSucceedAlert(message: string, timeout: number = 3000) {
    this.showAlert(message, 'success', timeout);
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

  /**
   * Show Or Hide Top Button
   * @param multiple Take the maximum percentage
   */
  showOrHideToTopBtn(multiple: number = 3) {
    const pos = (document.documentElement.scrollTop || document.body.scrollTop) + document.documentElement.offsetHeight;
    const max = document.documentElement.scrollHeight;
    if (pos > (max / multiple)) {
      this.topBtn = true;
    } else {
      this.topBtn = false;
    }
  }

  scrollToTop() {
    window.scrollTo(0, 0);
  }
}
