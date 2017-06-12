import { Component, OnInit } from '@angular/core';
import { ConnectionService } from "../../providers";
import { NgbModal, ModalDismissReasons } from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'app-connection',
  templateUrl: './connection.component.html',
  styleUrls: ['./connection.component.css']
})

export class ConnectionComponent implements OnInit {
  list: Array<string> = [];
  addUrl: string = '';
  constructor(private conn: ConnectionService, private modal: NgbModal) { }

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
        if(!this.addUrl) {
          return;
        }
        let data = new FormData();
        data.append('address', this.addUrl);
        this.conn.addConnection(data).subscribe(isOk => {
          if (isOk) {
            this.getAllConnections();
          }
        })
      }
    },err=>{});
  }
  remove(address: string) {
    let data = new FormData();
    data.append('address', address);
    this.conn.removeConnection(data).subscribe(isOk => {
      if (isOk) {
        this.getAllConnections();
      }
    })
  }
}