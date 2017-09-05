import { Injectable, ComponentRef, TemplateRef, ComponentFactory, ComponentFactoryResolver, Injector, ApplicationRef } from '@angular/core';
import { PopupWindow } from './popup-window';
import { PopupRef } from './popup-ref';
import { Router } from '@angular/router';

@Injectable()
export class PopupStack {
  private _windowFactory: ComponentFactory<PopupWindow>;
  constructor(private _componentFactoryResolver: ComponentFactoryResolver,
    private _injector: Injector, private _applicationRef: ApplicationRef,
    private router: Router) {
    this._windowFactory = _componentFactoryResolver.resolveComponentFactory(PopupWindow);
  }
  open(content: any, isAutoLeave: boolean = true) {
    let windowCmpRef: ComponentRef<PopupWindow>
    const containerEl = document.querySelector('body');
    let contentRef = null;
    if (!content || this.isString(content)) {
      console.log('content null');
      return;
    } else if (content instanceof TemplateRef) {
      contentRef = content.createEmbeddedView(content);
      windowCmpRef = this._windowFactory.create(this._injector, [contentRef.rootNodes]);
      this._applicationRef.attachView(contentRef);
      // viewRef.destroy();
    } else {
      const contentCmptFactory = this._componentFactoryResolver.resolveComponentFactory(content);
      contentRef = contentCmptFactory.create(this._injector);
      windowCmpRef = this._windowFactory.create(this._injector, [[contentRef.location.nativeElement]]);
      this._applicationRef.attachView(contentRef.hostView);
      // componentRef.destroy();
    }
    containerEl.appendChild(windowCmpRef.location.nativeElement);
    return { ref: new PopupRef().getRef(windowCmpRef, this.router, isAutoLeave), instance: contentRef.instance };
  }
  private isString(value: any): value is string {
    return typeof value === 'string';
  }
}
