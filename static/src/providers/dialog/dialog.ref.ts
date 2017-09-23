import { Injectable, ComponentRef } from '@angular/core';
import { DialogComponent } from './dialog.component';
import { DialogOverlayComponent } from './dialog.overlay.component';

@Injectable()
export class DialogRef {
  _resolve: (result?: any) => void;
  _reject: (reason?: any) => void;
  result: Promise<any>;
  constructor(private _contentRef: ComponentRef<DialogComponent>,
    private _overlaytRef?: ComponentRef<DialogOverlayComponent>) {
    this.result = new Promise((resolve, reject) => {
      this._resolve = resolve;
      this._reject = reject;
    });
    this.result.then(null, () => { });
  }

  close(result?: any) {
    this._resolve(result);
    this._removeModalElements();
  }
  private _removeModalElements() {
    const contentNativeEl = this._contentRef.location.nativeElement;
    contentNativeEl.parentNode.removeChild(contentNativeEl);
    this._contentRef.destroy();
    if (this._overlaytRef) {
      const overlayNativeEl = this._overlaytRef.location.nativeElement;
      overlayNativeEl.parentNode.removeChild(overlayNativeEl);
      this._overlaytRef.destroy();
      this._overlaytRef = null;
    }
    this._contentRef = null;
  }
}
