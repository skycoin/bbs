import { Component, OnInit, ViewEncapsulation, Output, EventEmitter } from '@angular/core';
import { ApiService, UserService } from "../../providers";
import { Board } from "../../providers/api/msg";
import { Router, ActivatedRoute } from "@angular/router";
import { NgbModal, ModalDismissReasons } from "@ng-bootstrap/ng-bootstrap";
import { FormControl, FormGroup } from '@angular/forms';

@Component({
    selector: 'boards-list',
    templateUrl: 'boards-list.component.html',
    styleUrls: ['boards.css'],
    encapsulation: ViewEncapsulation.None,
})
export class BoardsListComponent implements OnInit {
    @Output() board: EventEmitter<string> = new EventEmitter();
    boardsTitle: string = 'Boards List';
    errorMessage: string;
    isRoot: boolean = false;
    boards: Array<Board> = [];
    addForm = new FormGroup({
        name: new FormControl(),
        description: new FormControl(),
        seed: new FormControl()
    });
    constructor(private api: ApiService, private user: UserService, private router: Router, private modal: NgbModal) {

    }

    ngOnInit(): void {
        this.api.getBoards().subscribe(boards => {
            this.boards = boards;
        })
        this.api.getStats().subscribe(root => {
            this.isRoot = root;
        })
    }
    openAdd(content) {
        this.modal.open(content).result.then((result) => {
            if (result) {
                let data = new FormData();
                data.append('name', this.addForm.get('name').value);
                data.append('description', this.addForm.get('description').value);
                data.append('seed', this.addForm.get('seed').value);
                this.api.addBoard(data).subscribe(res => {
                    console.log('add board:', res);
                    this.api.getBoards().subscribe(boards => {
                        this.boards = boards;
                    });
                });
            }
        });
    }
    openThreads(ev: Event, key, url: string) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        this.router.navigate(['/threads', { board: key, url: url }])
        // this.board.emit(this.boards[0].public_key);
    }
    private getDismissReason(reason: any) {
        console.log('get dismiss reason:', reason);
    }
}
