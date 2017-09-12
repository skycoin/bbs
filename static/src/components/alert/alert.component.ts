import { Component, OnInit, ViewEncapsulation, HostBinding } from '@angular/core';
import { ActivePop } from '../../providers/popup/popup-stack';
import { bounceInAnimation } from '../../animations/common.animations';

@Component({
  selector: 'app-alert',
  templateUrl: './alert.component.html',
  styleUrls: ['./alert.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [bounceInAnimation]
})

export class AlertComponent implements OnInit {
  @HostBinding('style.display') display = 'block';
  @HostBinding('@bounceIn') animation = true;
  title = '';
  body = '';
  type = 'confirm';
  constructor(public activeModal: ActivePop) {
  }

  ngOnInit() {
  }
}
