import {
  Component,
  OnInit,
  Input,
  ViewEncapsulation,
  HostBinding,
  HostListener,
  Output,
  EventEmitter,
  ElementRef
} from '@angular/core';
import { PopupRef } from './popup-ref';

@Component({
  selector: 'popup-window',
  template: `<ng-content></ng-content>`,
  styleUrls: ['./pop-window.scss'],
  encapsulation: ViewEncapsulation.None,
  // tslint:disable-next-line:use-host-property-decorator
  host: {
    'tabindex': '-1',
    '[@fadeInOut]': ''
  },
})

// tslint:disable-next-line:component-class-suffix
export class PopupWindow implements OnInit {
  // tslint:disable-next-line:no-output-rename
  ref: PopupRef = null;
  canClick = true;
  constructor(private el: ElementRef) { }

  ngOnInit() { }
  @HostListener('click', ['$event'])
  _click(ev: Event) {
    if (this.canClick && this.ref && this.el.nativeElement === ev.target) {
      this.ref.close(false);
    }
  }
}
