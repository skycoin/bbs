import { Component, HostBinding, HostListener, OnInit, ViewEncapsulation } from '@angular/core';
import { CommonService, User, ApiService, Users, Alert, Popup } from '../../providers';
import { slideInLeftAnimation } from '../../animations/router.animations';
import { bounceInAnimation } from '../../animations/common.animations';
import { AlertComponent } from '../alert/alert.component';
import { FormControl, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-userlist',
  templateUrl: './userlist.component.html',
  styleUrls: ['./userlist.component.scss'],
  animations: [slideInLeftAnimation, bounceInAnimation],
  encapsulation: ViewEncapsulation.None
})
export class UserlistComponent implements OnInit {
  @HostBinding('@routeAnimation') routeAnimation = true;
  @HostBinding('style.display') display = 'block';
  userlist: Array<User> = [];
  editName = '';
  public addForm = new FormGroup({
    alias: new FormControl('', Validators.required),
    seed: new FormControl({ value: '', disabled: false }, Validators.required),
  });

  constructor(
    private api: ApiService,
    private common: CommonService,
    private alert: Alert,
    private pop: Popup) {
  }

  ngOnInit() {
    this.api.getAllUser().subscribe((users: Users) => {
      this.userlist = users.data.users;
    });
  }
  trackUsers(index, user) {
    return user ? user.alias : undefined;
  }
  openAdd(content: any) {
    this.addForm.reset();
    this.api.newSeed().subscribe(seed => {
      this.addForm.patchValue({ seed: seed.data })
      this.pop.open(content).result.then((result) => {
        if (result) {
          if (!this.addForm.valid) {
            this.alert.error({ content: 'Alias and Seed can not be empty' });
            return;
          }
          const data = new FormData();
          data.append('alias', this.addForm.get('alias').value);
          data.append('seed', this.addForm.get('seed').value);
          this.api.newUser(data).subscribe(res => {
            console.log('user:', res);
            if (res.okay) {
              this.userlist = res.data.users;
            }
          });
        }
      });
    })
  }

  getSeed(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.api.newSeed().subscribe(seed => {
      this.addForm.patchValue({ seed: seed.data })
    });
  }

  delUser(ev: Event, alias: string) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    if (!alias) {
      this.alert.error({ content: 'Parameter error!!!' });
      return;
    }
    const modalRef = this.pop.open(AlertComponent);
    modalRef.componentInstance.title = 'Delete User';
    modalRef.componentInstance.body = 'Do you delete the user?';
    modalRef.result.then(result => {
      if (result) {
        const data = new FormData();
        data.append('alias', alias);
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
