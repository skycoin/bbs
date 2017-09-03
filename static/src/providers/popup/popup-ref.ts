import { Injectable, ComponentRef } from '@angular/core';
import { PopupWindow } from './popup-window';
import { Router, NavigationStart } from '@angular/router';
import 'rxjs/add/operator/filter';

@Injectable()
export class PopupRef {
  _resolve: (result?: any) => void;
  _reject: (reason?: any) => void;
  _windowRef: ComponentRef<PopupWindow>
  result: Promise<any>;
  constructor() {
    this.result = new Promise((resolve, reject) => {
      this._resolve = resolve;
      this._reject = reject;
    });
    this.result.then(null, () => { });

  }
  getRef(ref: ComponentRef<PopupWindow>, router: Router, isAutoLeave = true) {
    if (router) {
      if (isAutoLeave) {
        router.events.filter(ev => ev instanceof NavigationStart).subscribe(() => {
          this.close();
        });
      }
    }
    this._windowRef = ref;
    return this;
  }
  close(result?: any) {
    if (this._windowRef) {
      this._resolve(result);
      this._removeModalElements();
    }
  }
  private _removeModalElements() {
    const windowNativeEl = this._windowRef.location.nativeElement;
    windowNativeEl.parentNode.removeChild(windowNativeEl);
    this._windowRef.destroy();
    this._windowRef = null;
  }
}
