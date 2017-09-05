import { Injectable, ComponentRef } from '@angular/core';
import { DialogComponent } from './dialog.component';
import { DialogOverlayComponent } from './dialog.overlay.component';

@Injectable()
export class DialogRef {

  constructor(private _contentRef: ComponentRef<DialogComponent>,
    private _overlaytRef?: ComponentRef<DialogOverlayComponent>) { }

  close() {
    if (this._contentRef) {
      this._removeModalElements();
    }
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
