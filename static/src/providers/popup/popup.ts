import { Injectable, TemplateRef, ComponentRef, ComponentFactory, ComponentFactoryResolver, Injector } from '@angular/core';
import { PopupStack, PopUpOptions } from './popup-stack';

@Injectable()
export class Popup {
  constructor(private stack: PopupStack) {
  }
  open(content: any, opts: PopUpOptions = { isAutoLeave: true, isDialog: true }) {
    return this.stack.open(content, opts);
  }
}
