import { Component, OnInit, HostBinding, HostListener } from '@angular/core';
import { UserService, User, CommonService } from "../../providers";
import { NgbModal } from "@ng-bootstrap/ng-bootstrap";
import { slideInLeftAnimation } from "../../animations/router.animations";
import { AlertComponent } from "../alert/alert.component";

@Component({
  selector: 'app-userlist',
  templateUrl: './userlist.component.html',
  styleUrls: ['./userlist.component.css'],
  animations: [slideInLeftAnimation],
})
export class UserlistComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  // @HostBinding('style.position') position = 'absolute';
  userlist: Array<User> = [];
  editName: string = '';
  constructor(private user: UserService, private modal: NgbModal, private common: CommonService) { }
  ngOnInit() {
    this.user.getAll().subscribe(userlist => {
      this.userlist = userlist;
    })
  }
  openEdit(content: any, key: string) {
    const modalRef = this.modal.open(content).result.then(result => {
      if (result) {
        this.edit(this.editName, key);
      }
    });
  }
  edit(name, key: string) {
    let data = new FormData();
    data.append('alias', name);
    data.append('user', key);
    this.user.newOrModifyUser(data).subscribe(res => {
      this.userlist = [];
      this.user.getAll().subscribe(userlist => {
        this.userlist = userlist;
        this.common.showAlert('successfully modified', 'success', 3000);
      })
    })
  }
  remove(ev: Event, key: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    let data = new FormData();
    data.append('user', key);
    const modalRef = this.modal.open(AlertComponent);
    modalRef.componentInstance.title = 'Delete User';
    modalRef.componentInstance.body = 'Do you delete the user?';
    modalRef.result.then(result => {
      if (result) {
        this.user.remove(data).subscribe(isOk => {
          if (isOk) {
            this.userlist = [];
            this.user.getAll().subscribe(userlist => {
              this.userlist = userlist;
              this.common.showAlert('successfully deleted', 'success', 1000);
            })
          } else {
            this.common.showAlert('failed to delete', 'success', 1000);
          }
        })
      }
    })

  }

  @HostListener('window:scroll', ['$event'])
  windowScroll(event) {
    this.common.showOrHideToTopBtn();
  }
}
