import { Injectable, ComponentRef } from '@angular/core';
import { PopupWindow } from './popup-window';
import { Router, NavigationStart } from '@angular/router';
import 'rxjs/add/operator/filter';

@Injectable()
export class PopupRef {
  private _resolve: (result?: any) => void;
  private _reject: (reason?: any) => void;
  result: Promise<any>;

  constructor(private _windowRef: ComponentRef<PopupWindow>, private router: Router, ) {
    this.result = new Promise((resolve, reject) => {
      this._resolve = resolve;
      this._reject = reject;
    });
    this.result.then(null, () => { });
    if (this.router) {
      this.router.events.filter(ev => ev instanceof NavigationStart).subscribe(() => {
        this.close();
      });
    }
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
