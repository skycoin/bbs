import { Component, OnInit, ViewEncapsulation, Input } from '@angular/core';
import { flyInOutAnimation, rotate45Animation } from './fab.animations';

@Component({
  selector: 'fab',
  templateUrl: './fab.components.html',
  styleUrls: ['./fab.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [flyInOutAnimation, rotate45Animation]
})

export class FabComponent implements OnInit {
  @Input() showAnimation = true;
  @Input() icon = 'fa-plus'
  isShow = false;
  fabAnimation = 'inactive';
  constructor() { }
  ngOnInit() {
  }
  open() {
    this.close();
  }
  close() {
    if (this.showAnimation) {
      this.isShow = !this.isShow;
      this.fabAnimation = this.fabAnimation === 'inactive' ? 'active' : 'inactive';
    }
  }
}
