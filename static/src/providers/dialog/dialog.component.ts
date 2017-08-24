import { Component, OnInit, ViewEncapsulation, ComponentRef, HostListener } from '@angular/core';
import { DialogAnimation } from './dialog.animation';

@Component({
  selector: 'app-dialog',
  templateUrl: './dialog.component.html',
  styleUrls: ['./dialog.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [DialogAnimation]
})

export class DialogComponent implements OnInit {
  self: ComponentRef<DialogComponent>;
  title = 'Test Title';
  body = 'Test Body';

  constructor() {
  }

  ngOnInit() {
  }
  @HostListener('click', ['$event'])
  _click(ev: Event) {
    const el = this.self.location.nativeElement;
    el.parentNode.removeChild(el);
    this.self.destroy();
    this.self = null;
  }
}
