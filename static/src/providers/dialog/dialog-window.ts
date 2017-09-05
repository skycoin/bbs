import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { DialogAnimation } from './dialog.animation';

@Component({
  selector: 'dialog-window',
  template: `<div role="dialog" class="dialog-container">
    <bbs-dialog></bbs-dialog>
  </div>`,
  styleUrls: ['./dialog-window.scss'],
  encapsulation: ViewEncapsulation.None
})

export class DialogWindowComponent implements OnInit {
  // alerts: Array<AlertData> = [];
  constructor() { }
  ngOnInit() { }

  hidden(index) {
    // this.alerts.splice(index, 1);
  }
}
