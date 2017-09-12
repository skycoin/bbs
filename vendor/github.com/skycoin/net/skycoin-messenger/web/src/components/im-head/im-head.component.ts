import { Component, OnInit, ViewEncapsulation, Input, HostListener } from '@angular/core';
import { HeadColorMatch, SocketService } from '../../providers';
import { MdDialog } from '@angular/material';
import { ImInfoDialogComponent } from '../im-info-dialog/im-info-dialog.component';

@Component({
  selector: 'app-im-head',
  templateUrl: './im-head.component.html',
  styleUrls: ['./im-head.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImHeadComponent implements OnInit {
  @Input() key = '';
  @Input() canClick = true;
  name = '';
  icon: HeadColorMatch = { bg: '#fff', color: '#000' };
  constructor(private socket: SocketService, private dialog: MdDialog) { }

  ngOnInit() {
    if (this.key === '') {
      this.name = '?';
    } else {
      const size = this.key.length;
      this.name = this.key.substr(size - 1, size).toUpperCase();
    }
    if (this.socket.userInfo.get(this.key) !== undefined) {
      const icon = this.socket.userInfo.get(this.key).Icon;
      if (icon !== undefined) {
        this.icon = icon;
      }
    }
  }

  @HostListener('click', ['$event'])
  _click(ev: Event) {
    if (!this.canClick) {
      return;
    }
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    const ref = this.dialog.open(ImInfoDialogComponent, {
      position: { top: '10%' },
      // panelClass: 'alert-dialog-panel',
      backdropClass: 'alert-backdrop',
      width: '23rem'
    });
    ref.componentInstance.key = this.key;
  }
}

