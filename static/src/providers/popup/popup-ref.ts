import { Injectable, ComponentRef } from '@angular/core';
import { PopupWindow } from './popup-window';
import { PopupBackdrop } from './popup.backdrop';
import { Router, NavigationStart } from '@angular/router';
import 'rxjs/add/operator/filter';

@Injectable()
export class PopupRef {
  _resolve: (result?: any) => void;
  _reject: (reason?: any) => void;
  _contentRef: ComponentRef<any>;
  _windowRef: ComponentRef<PopupWindow>
  _backdropRef: ComponentRef<PopupBackdrop>
  result: Promise<any>;
  constructor() {
    this.result = new Promise((resolve, reject) => {
      this._resolve = resolve;
      this._reject = reject;
    });
    this.result.then(null, () => { });

  }
  get componentInstance(): any {
    if (this._contentRef.instance) {
      return this._contentRef.instance;
    }
  }

  // only needed to keep TS1.8 compatibility
  set componentInstance(instance: any) { }
  getRef(ref: ComponentRef<PopupWindow>,
    backdropRef: ComponentRef<PopupBackdrop>, content: ComponentRef<any>, router: Router, isAutoLeave = true) {
    if (router) {
      if (isAutoLeave) {
        router.events.filter(ev => ev instanceof NavigationStart).subscribe(() => {
          this.close();
        });
      }
    }
    this._contentRef = content;
    this._windowRef = ref;
    this._backdropRef = backdropRef;
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
    if (this._backdropRef) {
      const backdropNativeEl = this._backdropRef.location.nativeElement;
      backdropNativeEl.parentNode.removeChild(backdropNativeEl);
      this._backdropRef.destroy();
    }
    this._windowRef.destroy();
    this._windowRef = null;
    this._backdropRef = null;
  }
}
