import { Injectable, ComponentFactoryResolver, Injector, ApplicationRef, ComponentRef } from '@angular/core';
import { LoadingComponent } from './loading';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/timer'

@Injectable()
export class LoadingService {
  private ref: ComponentRef<LoadingComponent>;
  constructor(private _componentFactoryResolver: ComponentFactoryResolver,
    private _injector: Injector, private _applicationRef: ApplicationRef) { }

  start(loadingText: string = 'Loading', options: LoadingOptions = { IconColor: '#000', TextColor: '#000' }) {
    if (!loadingText || loadingText === '') {
      loadingText = 'Loading';
    }
    Observable.timer(10).subscribe(() => {
      this.show(loadingText, options);
    });
    return Promise.resolve();
  }

  close() {
    this.ref.destroy();
    return Promise.resolve();
  }
  private show(loadingText: string, options?: LoadingOptions) {
    const containerEl = document.querySelector('body');
    const contentCmptFactory = this._componentFactoryResolver.resolveComponentFactory(LoadingComponent);
    this.ref = contentCmptFactory.create(this._injector);
    const instance = this.ref.instance;
    instance.loadingText = loadingText;
    instance.iconColor = options.IconColor;
    instance.textColor = options.TextColor;
    if (options.Duration && options.Duration > 0) {
      Observable.timer(options.Duration).subscribe(() => {
        this.close();
      })
    }
    this._applicationRef.attachView(this.ref.hostView);
    containerEl.appendChild(this.ref.location.nativeElement);
  }

}

export interface LoadingOptions {
  IconColor?: string;
  TextColor?: string;
  Duration?: number;
}
