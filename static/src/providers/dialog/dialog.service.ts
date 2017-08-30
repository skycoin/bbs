import { Injectable, ComponentFactoryResolver, Injector, ApplicationRef, ComponentRef, ComponentFactory } from '@angular/core';
import { DialogComponent } from './dialog.component';
import { DialogWindowComponent } from './dialog-window';
import { DialogRef } from './dialog.ref'
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/timer'

@Injectable()
export class Dialog {
  constructor(private _componentFactoryResolver: ComponentFactoryResolver,
    private _injector: Injector, private _applicationRef: ApplicationRef) {
  }

  open() {
    const containerEl = document.body;
    const windowCmptFactory = this._componentFactoryResolver.resolveComponentFactory(DialogWindowComponent);
    const ref = windowCmptFactory.create(this._injector);
    this._applicationRef.attachView(ref.hostView);
    containerEl.appendChild(ref.location.nativeElement);
  }
}
