import { Injectable, ComponentFactoryResolver, Injector, ApplicationRef, ComponentRef, ComponentFactory } from '@angular/core';
import { DialogComponent } from './dialog.component';
import { DialogRef } from './dialog.ref'
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/timer'

@Injectable()
export class Dialog {
  constructor(private _componentFactoryResolver: ComponentFactoryResolver,
    private _injector: Injector, private _applicationRef: ApplicationRef) {
  }

  open() {
    Observable.timer(10).subscribe(() => {
      const containerEl = document.querySelector('body');
      const contentCmptFactory = this._componentFactoryResolver.resolveComponentFactory(DialogComponent);
      const ref = contentCmptFactory.create(this._injector);
      ref.instance.self = ref;
      this._applicationRef.attachView(ref.hostView);
      containerEl.appendChild(ref.location.nativeElement);
    });
  }
}
