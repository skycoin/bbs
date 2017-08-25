import { Component, HostBinding, HostListener, OnInit } from '@angular/core';
import { CommonService, User, ApiService, Users, Alert } from '../../providers';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { AlertComponent } from '../alert/alert.component';
import { FormControl, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-userlist',
  templateUrl: './userlist.component.html',
  styleUrls: ['./userlist.component.scss'],
  animations: [slideInLeftAnimation],
})
export class UserlistComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  userlist: Array<User> = [];
  editName = '';
  public addForm = new FormGroup({
    alias: new FormControl('', Validators.required),
    seed: new FormControl('', Validators.required),
  });

  constructor(
    private api: ApiService,
    private modal: NgbModal,
    private common: CommonService,
    private alert: Alert) {
  }

  ngOnInit() {
    this.api.getAllUser().subscribe((users: Users) => {
      this.userlist = users.data.users;
    });
  }

  // openEdit(content: any, key: string) {
  //   this.editName = '';
  //   if (key === '') {
  //     // this.common.showErrorAlert('Parameter error!!!');
  //     return;
  //   }
  //   const modalRef = this.modal.open(content).result.then(result => {
  //     if (result) {
  //       this.edit(this.editName, key);
  //     }
  //   });
  // }

  openAdd(content: any) {
    this.addForm.reset();
    this.modal.open(content).result.then((result) => {
      if (result) {
        if (!this.addForm.valid) {
          // this.common.showErrorAlert('Alias and Seed can not be empty');
          return;
        }
        const data = new FormData();
        data.append('alias', this.addForm.get('alias').value);
        data.append('seed', this.addForm.get('seed').value);
        // this.user.newMaster(data).subscribe(user => {
        //   this.userlist.unshift(user);
        // });
      }
    });
  }

  // edit(name, key: string) {
  //   if (name === '') {
  //     // this.common.showErrorAlert('UserName can not be empty');
  //     return;
  //   }
  //   const data = new FormData();
  //   data.append('alias', name);
  //   data.append('user', key);
  //   this.user.newOrModifyUser(data).subscribe(res => {
  //     this.userlist = [];
  //     this.user.getAll().subscribe(userlist => {
  //       this.userlist = userlist;
  //       // this.common.showAlert('modified successfully', 'success', 3000);
  //     });
  //   });
  // }

  delUser(ev: Event, key: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    const modalRef = this.modal.open(AlertComponent);
    modalRef.componentInstance.title = 'Delete User';
    modalRef.componentInstance.body = 'Do you delete the user?';
    modalRef.result.then(result => {
      if (result) {
        const data = new FormData();
        data.append('alias', key);
        this.api.delUser(data).subscribe((users: Users) => {
          if (users.okay) {
            this.userlist = users.data.users;
            this.alert.success({ content: 'deleted successfully' });
          }
        })
      }
    }, err => { });

  }
}
