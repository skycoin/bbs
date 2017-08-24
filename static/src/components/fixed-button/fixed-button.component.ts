import { Component, ViewEncapsulation, HostListener, OnDestroy, HostBinding } from '@angular/core';

@Component({
  selector: 'app-fixed-button',
  templateUrl: './fixed-button.component.html',
  styleUrls: ['./fixed-button.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class FixedButtonComponent implements OnDestroy {
  @HostBinding('style.display') display = 'none';
  handle: Function = null;
  constructor() { }
  ngOnDestroy() {
    this.handle = null;
  }
  @HostListener('click', ['$event'])
  _click(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    if (this.handle !== null) {
      this.handle();
    }
  }
}
