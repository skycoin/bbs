import { Component, HostBinding, OnInit, ViewEncapsulation } from '@angular/core';
import { CommonService, ApiService, Connnections, Connnection } from '../../providers';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { AlertComponent } from '../../components/alert/alert.component';
@Component({
  selector: 'app-connection',
  templateUrl: './connection.component.html',
  styleUrls: ['./connection.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [slideInLeftAnimation],
})

export class ConnectionComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  list: Array<Connnection> = [];
  addUrl = '';

  constructor(
    private modal: NgbModal,
    private api: ApiService) {
  }

  ngOnInit() {
    this.api.getAllConnections().subscribe((conns: Connnections) => {
      console.log('init conns1:', conns);
      this.list = conns.data.connections;
      console.log('init conns2:', this.list);
    });
  }

  openAdd(content) {
    this.addUrl = '';
    this.modal.open(content).result.then((result) => {
      if (result) {
        if (!this.addUrl) {
          // this.common.showAlert('The link can not be empty', 'danger', 3000);
          return;
        }
        const data = new FormData();
        data.append('address', this.addUrl);
        this.api.newConnection(data).subscribe((conns: Connnections) => {
          if (conns.okay) {
            this.list = conns.data.connections;
          }
        })
      }
    }, err => {
    });
  }

  remove(address: string) {
    const modalRef = this.modal.open(AlertComponent);
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
