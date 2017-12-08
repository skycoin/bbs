import {
  Injectable,
  ComponentRef,
  TemplateRef,
  ComponentFactory,
  ComponentFactoryResolver,
  Injector,
  ApplicationRef,
  ReflectiveInjector,
} from '@angular/core';
import { PopupWindow } from './popup-window';
import { PopupBackdrop } from './popup.backdrop'
import { PopupRef } from './popup-ref';
import { Router } from '@angular/router';

@Injectable()
export class ActivePop {
  close(result?: any): void { }
}

@Injectable()
export class PopupStack {
  private _windowFactory: ComponentFactory<PopupWindow>;
  private _backdropFactory: ComponentFactory<PopupBackdrop>;
  constructor(private _componentFactoryResolver: ComponentFactoryResolver,
    private _injector: Injector, private _applicationRef: ApplicationRef,
    private router: Router) {
    this._windowFactory = _componentFactoryResolver.resolveComponentFactory(PopupWindow);
    this._backdropFactory = _componentFactoryResolver.resolveComponentFactory(PopupBackdrop);
  }
  open(content: any, opts: PopUpOptions = { isAutoLeave: true, isDialog: true, canClickBackdrop: true }) {
    let windowCmpRef: ComponentRef<PopupWindow>
    const containerEl = document.body;
    const activeModal = new ActivePop();
    let contentRef = null;
    let backdropCmptRef = null;
    if (!content || this.isString(content)) {
      console.log('content null');
      return;
    } else if (content instanceof TemplateRef) {
      contentRef = content.createEmbeddedView(activeModal);
      windowCmpRef = this._windowFactory.create(this._injector, [contentRef.rootNodes]);
      this._applicationRef.attachView(contentRef);
    } else {
      const contentCmptFactory = this._componentFactoryResolver.resolveComponentFactory(content);
      const modalContentInjector =
        ReflectiveInjector.resolveAndCreate([{ provide: ActivePop, useValue: activeModal }], this._injector);
      contentRef = contentCmptFactory.create(modalContentInjector);
      windowCmpRef = this._windowFactory.create(this._injector, [[contentRef.location.nativeElement]]);
      this._applicationRef.attachView(contentRef.hostView);
    }
    const ref = new PopupRef().getRef(windowCmpRef, backdropCmptRef, contentRef, this.router, opts.isAutoLeave);
    if (opts.isDialog) {
      windowCmpRef.location.nativeElement.classList.add('popup-dialog');
      backdropCmptRef = this._backdropFactory.create(this._injector);
      this._applicationRef.attachView(backdropCmptRef.hostView);
      containerEl.appendChild(backdropCmptRef.location.nativeElement);
      ref._backdropRef = backdropCmptRef;
      windowCmpRef.instance.ref = ref;
      windowCmpRef.instance.canClick = opts.canClickBackdrop;
    }
    activeModal.close = (result: any) => { ref.close(result); };
    containerEl.appendChild(windowCmpRef.location.nativeElement);
    return ref;
  }
  private isString(value: any): value is string {
    return typeof value === 'string';
  }
}

export interface PopUpOptions {
  isDialog?: boolean;
  isAutoLeave?: boolean;
  canClickBackdrop?: boolean;
}



