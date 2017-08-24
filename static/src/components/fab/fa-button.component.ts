import { Component, Input } from '@angular/core';

@Component({
  selector: 'fa-button',
  template: `<a href="javascript:void(0);">
    <div class="fab-tooltip">{{name}}</div>
    <i class="fa" [ngClass]="icon" aria-hidden="true"></i>
    </a>
  `,
  styleUrls: ['./fab.component.scss']
})

// tslint:disable-next-line:component-class-suffix
export class FabButton {
  @Input() name = '';
  @Input() icon = '';
  constructor() { }
}
