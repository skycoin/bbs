import { Component, OnInit, ViewEncapsulation, HostBinding } from '@angular/core';
import { ApiService, UserService, CommonService } from '../../providers';
import { Board, UIOptions } from '../../providers/api/msg';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, ModalDismissReasons } from '@ng-bootstrap/ng-bootstrap';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { AlertComponent } from '../alert/alert.component';
import { slideInLeftAnimation } from '../../animations/router.animations';

@Component({
    selector: 'app-boardslist',
    templateUrl: 'boards-list.component.html',
    styleUrls: ['boards-list.scss'],
    encapsulation: ViewEncapsulation.None,
    animations: [slideInLeftAnimation],
})
export class BoardsListComponent implements OnInit {
    @HostBinding('@routeAnimation') routeAnimation = true;
    @HostBinding('style.display') display = 'block';
    public sort = 'desc';
    public isRoot = false;
    public boards: Array<Board> = [];
    public subscribeForm = new FormGroup({
        address: new FormControl('', Validators.required),
        board: new FormControl('', Validators.required)
    });
    public addForm = new FormGroup({
        name: new FormControl(),
        description: new FormControl(),
        seed: new FormControl()
    });
    public tmpBoard: Board = null;
    constructor(
        private api: ApiService,
        private user: UserService,
        private router: Router,
        private modal: NgbModal,
        public common: CommonService) {
    }

    ngOnInit(): void {
        this.common.loading.start();
        this.getBoards();
        this.api.getStats().subscribe(root => {
            this.isRoot = root;
        });
    }
    setSort() {
        this.sort = this.sort === 'desc' ? 'esc' : 'desc';
    }
    getBoards() {
        this.api.getBoards().subscribe(boards => {
            if (!boards || boards.length <= 0) {
                this.common.loading.close();
                return;
            }
            this.boards = boards;
            this.boards.forEach(el => {
                if (!el || !el.public_key) {
                    return;
                }
                const data = new FormData();
                data.append('board', el.public_key);
                this.api.getSubscription(data).subscribe(res => {
                    el.ui_options = { subscribe: true };
                })
                this.common.loading.close();
            });
        });
    }
    openInfo(ev: Event, board: Board, content: any) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        if (!board) {
            this.common.showErrorAlert('Failed to get info!!')
            return;
        }
        this.tmpBoard = board;
        this.modal.open(content, { size: 'lg' });
    }
    openAdd(content) {
        this.addForm.reset();
        this.modal.open(content).result.then((result) => {
            if (result === true) {
                const data = new FormData();
                data.append('name', this.addForm.get('name').value);
                data.append('description', this.addForm.get('description').value);
                data.append('seed', this.addForm.get('seed').value);
                this.api.addBoard(data).subscribe(res => {
                    this.getBoards();
                    this.common.showSucceedAlert('Added Successfully');
                });
            }
        }, err => { });
    }
    subscribe(ev: Event, content: any) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        this.modal.open(content).result.then(result => {
            if (result) {
                if (!this.subscribeForm.valid) {
                    this.common.showErrorAlert('The Board Key Or Address can not be empty!!!');
                    return;
                }
                const data = new FormData()
                data.append('address', this.subscribeForm.get('address').value);
                data.append('board', this.subscribeForm.get('board').value)
                this.api.subscribe(data).subscribe(isOk => {
                    if (isOk) {
                        this.common.showSucceedAlert('Subscribe successfully');
                        this.getBoards();
                    }
                })
            }
        })

    }
    unSubscribe(ev: Event, boardKey: string) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        if (boardKey === '') {
            this.common.showErrorAlert('UnSubscribe failed');
            return;
        }
        const data = new FormData();
        data.append('board', boardKey);
        this.api.unSubscribe(data).subscribe(isOk => {
            if (isOk) {
                this.common.showSucceedAlert('Unsubscribe successfully');
                this.getBoards();
            }
        })
    }
    openThreads(ev: Event, key: string) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        if (!key) {
            this.common.showErrorAlert('Abnormal parameters!!!', 3000);
            return;
        }
        this.router.navigate(['/threads', { board: key }])
    }
}
