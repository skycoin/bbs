import { Http, Response } from '@angular/http';
import { Injectable } from '@angular/core';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/filter';
import { Observable } from 'rxjs/Observable';
import { LoadingComponent, FixedButtonComponent } from '../../components';

@Injectable()
export class CommonService {
  public alertType = 'info';
  public alertMessage = 'test alert';
  public alert = false;
  public topBtn = false;
  public fb: FixedButtonComponent = null;
  public loading: LoadingComponent = null;
  // public sortBy = 'desc';
  constructor(private http: Http) {
  }

  handleError(error: Response) {
    if (this.loading) {
      this.loading.close();
    }
    console.error('Error:', error.json() || 'Server error', 'danger');
    this.showAlert((error.json() instanceof Object ? 'Server error' : error.json()) || 'Server error', 'danger', 3000);
    return Observable.throw(error.json() || 'Server error');
  }

  handleGet(url: string) {
    if (!url) {
      return Observable.throw('The connection is empty');
    }
    return this.http.get(url)
      .filter((res: Response) => res.status === 200)
      .map((res: Response) => res.json()).catch(err => this.handleError(err));
  }

  handlePost(url: string, data: FormData) {
    if (!url || !data) {
      return Observable.throw('Parameters and connections can not be empty');
    }
    return this.http.post(url, data)
      .filter((res: Response) => res.status === 200)
      .map((res: Response) => res.json()).catch(err => this.handleError(err));
  }

  copy(ev) {
    if (ev) {
      this.showSucceedAlert('Copy Successful');
    } else {
      this.showErrorAlert('Copy Failed');
    }
  }

  /**
   * Show Error Alert
   * @param message Error Text
   * @param timeout
   */
  showErrorAlert(message: string, timeout: number = 3000) {
    this.showAlert(message, 'danger', timeout);
  }
  showWarningAlert(message: string, timeout: number = 3000) {
    this.showAlert(message, 'warning', timeout);
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
  showOrHideToTopBtn() {
    const pos = (document.documentElement.scrollTop || document.body.scrollTop) + document.documentElement.offsetHeight;
    const max = document.documentElement.scrollHeight;
    const clientHeight = document.documentElement.clientHeight;
    if (pos > max - (max - clientHeight)) {
      this.topBtn = true;
    } else if (pos <= clientHeight) {
      this.topBtn = false;
    }
  }

  scrollToTop() {
    window.scrollTo(0, 0);
  }
}
