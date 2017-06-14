import { Component, OnInit, ViewEncapsulation, Output, EventEmitter,HostBinding } from '@angular/core';
import { ApiService, UserService, CommonService } from "../../providers";
import { Board, UIOptions } from "../../providers/api/msg";
import { Router, ActivatedRoute } from "@angular/router";
import { NgbModal, ModalDismissReasons } from "@ng-bootstrap/ng-bootstrap";
import { FormControl, FormGroup } from '@angular/forms';
import { AlertComponent } from "../alert/alert.component";
import { slideInLeftAnimation } from "../../animations/router.animations";

@Component({
    selector: 'boards-list',
    templateUrl: 'boards-list.component.html',
    styleUrls: ['boards.css'],
    encapsulation: ViewEncapsulation.None,
    animations: [slideInLeftAnimation],
})
export class BoardsListComponent implements OnInit {
    @HostBinding('@routeAnimation') routeAnimation = true;
    @HostBinding('style.display') display = 'block';
    @HostBinding('style.position') position = 'absolute';
    @Output() board: EventEmitter<string> = new EventEmitter();
    private isRoot: boolean = false;
    private boards: Array<Board> = [];
    private addForm = new FormGroup({
        name: new FormControl(),
        description: new FormControl(),
        seed: new FormControl()
    });
    private tmpBoard: Board = null;
    constructor(
        private api: ApiService,
        private user: UserService,
        private router: Router,
        private modal: NgbModal,
        private common: CommonService) {

    }

    ngOnInit(): void {
        this.getBoards();
        this.api.getStats().subscribe(root => {
            this.isRoot = root;
        });
    }

    private getBoards() {
        this.api.getBoards().subscribe(boards => {
            this.boards = boards;
            this.boards.forEach(el => {
                let data = new FormData();
                data.append('board', el.public_key);
                this.api.getSubscription(data).subscribe(res => {
                    if (res.config && res.config.secret_key) {
                        el.ui_options = { subscribe: true };
                    }
                })
            });
        });
    }

    openInfo(ev: Event, board: Board, content: any) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        this.tmpBoard = board;
        this.modal.open(content, { size: 'lg' });
    }
    openAdd(content) {
        this.modal.open(content).result.then((result) => {
            if (result === true) {
                let data = new FormData();
                data.append('name', this.addForm.get('name').value);
                data.append('description', this.addForm.get('description').value);
                data.append('seed', this.addForm.get('seed').value);
                this.api.addBoard(data).subscribe(res => {
                    this.api.getBoards().subscribe(boards => {
                        this.boards = boards;
                        this.common.showAlert('Added successfully', 'success', 3000);
                    });
                });
            }
        }, err => { });
    }
    subscribe(ev: Event, key: string, index: number) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        let data = new FormData()
        data.append('board', key)
        if (!this.boards[index].ui_options.subscribe) {
            this.api.subscribe(data).subscribe(isOk => {
                let options = { subscribe: isOk };
                this.boards[index].ui_options = options;
                this.common.showAlert('Subscribe successfully', 'success', 3000);
            })
        } else {
            this.api.unSubscribe(data).subscribe(isOk => {
                if (isOk) {
                    this.boards[index].ui_options.subscribe = false;
                    this.common.showAlert('Unsubscribe successfully', 'success', 3000);
                    this.getBoards();
                }
            })
        }
    }
    openThreads(ev: Event, key, url: string) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        this.router.navigate(['/threads', { board: key }])
        // this.board.emit(this.boards[0].public_key);
    }
    private getDismissReason(reason: any) {
        console.log('get dismiss reason:', reason);
    }
}
