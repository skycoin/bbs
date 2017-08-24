import { Component, OnInit, Input, ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'popup-window',
  template: `<div role="button"><ng-content></ng-content></div>`,
  encapsulation: ViewEncapsulation.None,
  // tslint:disable-next-line:use-host-property-decorator
  host: {
    'role': 'dialog',
    'tabindex': '-1',
    '[@fadeInOut]': ''
  },
})

// tslint:disable-next-line:component-class-suffix
export class PopupWindow implements OnInit {
  @Input() windowClass: string;
  constructor() { }

  ngOnInit() { }
}
