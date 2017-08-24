import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { flyInOutAnimation, rotate45Animation } from './fab.animations';

@Component({
  selector: 'fab',
  templateUrl: './fab.components.html',
  styleUrls: ['./fab.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [flyInOutAnimation, rotate45Animation]
})

export class FabComponent implements OnInit {
  isShow = false;
  fabAnimation = 'inactive';
  constructor() { }
  ngOnInit() {
  }
  open(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.close();
  }
  close() {
    this.isShow = !this.isShow;
    this.fabAnimation = this.fabAnimation === 'inactive' ? 'active' : 'inactive';
  }
}
