import { Component, OnInit } from '@angular/core';
import { ConnectionService, CommonService } from "../../providers";
import { NgbModal, ModalDismissReasons } from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'app-connection',
  templateUrl: './connection.component.html',
  styleUrls: ['./connection.component.css']
})

export class ConnectionComponent implements OnInit {
  list: Array<string> = [];
  addUrl: string = '';
  constructor(
    private conn: ConnectionService,
    private modal: NgbModal,
    private common: CommonService) { }

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
        let data = new FormData();
        data.append('address', this.addUrl);
        this.conn.addConnection(data).subscribe(isOk => {
          if (isOk) {
            this.getAllConnections();
            this.common.showAlert('The connection was added successfully', 'success', 3000);
          }
        })
      }
    }, err => { });
  }
  remove(address: string) {
    let data = new FormData();
    data.append('address', address);
    this.conn.removeConnection(data).subscribe(isOk => {
      if (isOk) {
        this.getAllConnections();
        this.common.showAlert('The connection has been deleted', 'success', 3000);
      }
    })
  }
}