import { Component, OnInit } from '@angular/core';
import { UserService, User } from "../../providers";
import { NgbModal } from "@ng-bootstrap/ng-bootstrap";
import { ModalComponent } from "../modal/modal.component";

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {
  private user: User = null;
  closeResult: string;
  constructor(private userService: UserService, private modal: NgbModal) { }

  ngOnInit() {
    this.userService.getCurrent().subscribe(user => {
      this.user = user;
    })
  }

  editUserAlias(content) {
    const modalRef = this.modal.open(ModalComponent);
    modalRef.result.then(result => {
      if (result.ok) {
        this.edit(result.name);
      }
    }, err => {

    })
  }
  edit(name: string) {
    let data = new FormData();
    data.append('alias', name);
    data.append('user', this.user.public_key);
    this.userService.newOrModifyUser(data).subscribe(res => {
      console.log('edit :', res);
    })
  }
}
