import { Component, OnInit, ViewEncapsulation, Input, HostListener, Output, EventEmitter, HostBinding } from '@angular/core';
import { AlertData } from './msg';
import { AlertAnimation } from './alert.animation';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/timer'
@Component({
  selector: 'alert',
  templateUrl: './alert.html',
  styleUrls: ['./alert.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [AlertAnimation]
})

export class AlertComponent implements OnInit {
  @HostBinding('style.display') display = 'block';
  @HostBinding('@alertInOut') animation = true;
  @Input() data: AlertData;
  @Input() index;
  @Output() hidden = new EventEmitter<number>();
  constructor() { }

  ngOnInit() {
    if (this.data.autoDismiss && this.data.duration && this.data.duration > 0) {
      Observable.timer(this.data.duration).subscribe(() => {
        this.close();
      })
    }
  }
  getIcon(type: string) {
    let icon = '';
    switch (type) {
      case 'success':
        icon = 'fa-check-circle'
        break;
      case 'warning':
        icon = 'fa-exclamation-triangle';
        break;
      case 'error':
        icon = 'fa-ban';
        break;
      case 'info':
        icon = 'fa-info-circle';
        break;
      default:
        icon = 'fa-info-circle'
        break;
    }
    return icon;
  }
  close() {
    if (this.data.clickEvent) {
      this.data.clickEvent();
    }
    this.hidden.emit(this.index);
  }
  @HostListener('click', ['$event'])
  _click(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.close();
  }
}
