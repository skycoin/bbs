import { Injectable, ComponentFactoryResolver, Injector, ApplicationRef, ComponentRef } from '@angular/core';
import { AlertWindowComponent } from './alert-window';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/timer'
import { AlertData } from './msg';

@Injectable()
export class Alert {
  static Ref: ComponentRef<AlertWindowComponent>;
  constructor(private _componentFactoryResolver: ComponentFactoryResolver,
    private _injector: Injector, private _applicationRef: ApplicationRef) { }

  success(data?: AlertData) {
    if (!data.title) {
      data.title = 'Success';
    }
    data.type = 'success';
    this.start(data);
  }
  info(data?: AlertData) {
    if (!data.title) {
      data.title = 'Info';
    }
    data.type = 'info';
    this.start(data);
  }
  warning(data?: AlertData) {
    if (!data.title) {
      data.title = 'Warning';
    }
    data.type = 'warning';
    this.start(data);
  }
  error(data?: AlertData) {
    if (!data.title) {
      data.title = 'Error';
    }
    data.type = 'error';
    this.start(data);
  }
  private start(data?: AlertData) {
    if (!data.footer) {
      data.footer = 'Click to Dismiss';
    }
    if (data.autoDismiss === undefined) {
      data.autoDismiss = true;
    }
    if (data.autoDismiss) {
      data.duration = 3000;
    }
    if (Alert.Ref) {
      Alert.Ref.instance.alerts.unshift(data);
      return;
    }
    Observable.timer(10).subscribe(() => {
      this.create(data);
    });
  }
  create(data?: AlertData) {
    if (document.querySelector('alert-window')) {
      this.start(data);
      return;
    }
    const containerEl = document.body;
    const contentCmptFactory = this._componentFactoryResolver.resolveComponentFactory(AlertWindowComponent);
    Alert.Ref = contentCmptFactory.create(this._injector);
    Alert.Ref.instance.alerts.push(data);
    this._applicationRef.attachView(Alert.Ref.hostView);
    containerEl.appendChild(Alert.Ref.location.nativeElement);
  }
}
