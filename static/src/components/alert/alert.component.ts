import { Component, OnInit, ViewEncapsulation, HostBinding } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
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

  constructor(public activeModal: NgbActiveModal) {
  }

  ngOnInit() {
  }
}
