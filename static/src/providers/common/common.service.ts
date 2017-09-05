import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import 'rxjs/add/observable/throw';
import { Observable } from 'rxjs/Observable';
import { FixedButtonComponent } from '../../components';
import { Alert } from '../alert/alert.service';
import { LoadingService } from '../loading/loading.service';

@Injectable()
export class CommonService {
  public fb: FixedButtonComponent = null;
  constructor(private http: HttpClient, private alert: Alert, private loading: LoadingService) {
  }

  replaceURL(str: string) {
    let start = -1;
    let end = -1;
    start = str.indexOf('https://');
    if (start <= -1) {
      start = str.indexOf('http://');
    }
    if (start <= -1) {
      return str;
    }
    end = str.indexOf('.', start + 1);
    end = str.indexOf(' ', end + 1)
    const result = str.substring(start, end);
    return str.substring(0, start) + `<a href="${result}">${result}</a>` + str.substring(end);
  }

  replaceHtmlEnter(str: string) {
    return this.replaceStr(str, '\n', ' <br>');
  }
  replaceAt(str: string, start, end: number, replacement: string) {
    return str.substr(0, start) + replacement + str.substr(end);
  }
  replaceStr(str, seachStr, replaceStr: string) {
    if (str.length <= 0) {
      return str;
    }
    let pos = str.indexOf(seachStr);
    while (pos !== -1) {
      str = this.replaceAt(str, pos, pos + 1, replaceStr);
      pos = str.indexOf(seachStr, pos + 1);
    }
    return str;
  }
  stringToHex(str) {
    const arr = [];
    for (let i = 0; i < str.length; i++) {
      arr[i] = ('00' + str.charCodeAt(i).toString(16)).slice(-4);
    }
    return '\\u' + arr.join('\\u');
  }
  handleError(errorResponse: HttpErrorResponse) {
    const json = errorResponse.error ? errorResponse.error.error : undefined;
    this.alert.error(
      {
        title: json ? json.title : errorResponse.statusText,
        content: json ? json.details : errorResponse.message
      });
    this.loading.close();
    return Observable.throw(json ? json.details : errorResponse.message);
  }

  handleGet(url: string) {
    if (!url) {
      return Observable.throw('The connection is empty');
    }
    return this.http.get(url)
      .catch(err => this.handleError(err));
  }

  handlePost(url: string, data: FormData) {
    if (!url || !data) {
      return Observable.throw('Parameters and connections can not be empty');
    }
    return this.http.post(url, data)
      .catch(err => this.handleError(err));
  }

}
