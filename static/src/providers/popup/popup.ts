import { Injectable, TemplateRef, ComponentRef, ComponentFactory, ComponentFactoryResolver, Injector } from '@angular/core';
import { PopupStack } from './popup-stack';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/timer'

@Injectable()
export class Popup {
  constructor(private stack: PopupStack) {
  }
  open(content: any) {
    Observable.timer(10).subscribe(() => {
      return this.stack.open(content);
    });
  }
}
