import { Component, OnInit, ViewEncapsulation, ComponentRef, HostListener, HostBinding } from '@angular/core';
import { DialogAnimation } from './dialog.animation';

@Component({
  selector: 'bbs-dialog',
  templateUrl: './dialog.component.html',
  styleUrls: ['./dialog.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [DialogAnimation]
})

export class DialogComponent implements OnInit {
  @HostBinding('@dialogInOut') animation = true;
  @HostBinding('style.display') display = 'block';
  title = 'Test Title';
  body = 'Test Body';

  constructor() {
  }

  ngOnInit() {
  }
  @HostListener('click', ['$event'])
  _click(ev: Event) {
    console.log('click dialog');
  }
}
