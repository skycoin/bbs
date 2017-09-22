import { Component, HostBinding, OnInit, ViewEncapsulation } from '@angular/core';
import { CommonService, ApiService, Connnections, Connnection, Popup, Alert } from '../../providers';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { bounceInAnimation } from '../../animations/common.animations';
import { AlertComponent } from '../../components/alert/alert.component';
@Component({
  selector: 'app-connection',
  templateUrl: './connection.component.html',
  styleUrls: ['./connection.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation, bounceInAnimation],
})

export class ConnectionComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  list: Array<Connnection> = [];
  addUrl = '';

  constructor(
    private api: ApiService,
    private pop: Popup,
    private alert: Alert) {
  }

  ngOnInit() {
    this.api.getAllConnections().subscribe((conns: Connnections) => {
      this.list = conns.data.connections;
    });
  }

  openAdd(content) {
    this.addUrl = '';
    this.pop.open(content).result.then((result) => {
      if (result) {
        if (!this.addUrl) {
          this.alert.error({ content: 'The link can not be empty' });
          return;
        }
        const data = new FormData();
        data.append('address', this.addUrl);
        this.api.newConnection(data).subscribe((conns: Connnections) => {
          if (conns.okay) {
            this.list = conns.data.connections;
            this.alert.success({ content: 'Added Successfully' });
          }
        })
      }
    }, err => {
    });
  }

  remove(address: string) {
    const modalRef = this.pop.open(AlertComponent);
    modalRef.componentInstance.title = 'Delete Connection';
    modalRef.componentInstance.body = 'Do you delete the connection?';
    modalRef.result.then(result => {
      if (result) {
        const data = new FormData();
        data.append('address', address);
        this.api.delConnection(data).subscribe((conns: Connnections) => {
          if (conns.okay) {
            this.list = conns.data.connections;
          }
        })
      }
    }, err => { });
  }
}
