import { Component, HostBinding, OnInit, ViewEncapsulation } from '@angular/core';
import { CommonService, ConnectionService } from '../../providers';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { AlertComponent } from '../../components';
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
  list: Array<string> = [];
  addUrl = '';

  constructor(private conn: ConnectionService,
    private modal: NgbModal,
    private common: CommonService) {
  }

  ngOnInit() {
    this.getAllConnections();
  }

  getAllConnections() {
    this.conn.getAllConnections().subscribe(list => {
      this.list = list;
    });
  }

  openAdd(content) {
    this.addUrl = '';
    this.modal.open(content).result.then((result) => {
      if (result) {
        if (!this.addUrl) {
          this.common.showAlert('The link can not be empty', 'danger', 3000);
          return;
        }
        const data = new FormData();
        data.append('address', this.addUrl);
        this.conn.addConnection(data).subscribe(isOk => {
          if (isOk) {
            this.getAllConnections();
            this.common.showAlert('The connection was added successfully', 'success', 3000);
          }
        });
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
        this.conn.removeConnection(data).subscribe(isOk => {
          if (isOk) {
            this.getAllConnections();
            this.common.showAlert('The connection has been deleted', 'success', 3000);
          }
        });
      }
    }, err => { });
  }
}
