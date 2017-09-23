import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { AlertData } from './msg';
@Component({
  selector: 'alert-window',
  template: `<div role="alert">
    <alert (hidden)="hidden($event)" *ngFor="let alert of alerts;let i =index" [data]="alert" [index]="i"></alert>
  </div>`,
  styles: [`alert-window {
    display: block;
    position: fixed;
    top: 0;
    right: 0;
    width: 25%;
    padding: .5rem;
    z-index: 2010;
  }
  `],
  encapsulation: ViewEncapsulation.None,
})

export class AlertWindowComponent implements OnInit {
  alerts: Array<AlertData> = [];
  constructor() { }
  ngOnInit() { }

  hidden(index) {
    this.alerts.splice(index, 1);
  }
}
