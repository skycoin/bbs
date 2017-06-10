import { Component, OnInit, ViewEncapsulation, Output, EventEmitter } from '@angular/core';
import { ApiService, UserService } from "../../providers";
import { Board } from "../../providers/api/msg";
import { Router, ActivatedRoute } from "@angular/router";

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
    test: false;
    boards: Array<Board> = [];
    constructor(private api: ApiService, private user: UserService, private router: Router) {

    }

    ngOnInit(): void {
        this.api.getBoards().subscribe(boards => {
            this.boards = boards;
        })
        this.user.getCurrent().subscribe(user => {
            console.log('user', user);
        });
    }
    openThreads(ev: Event, key, url: string) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        this.router.navigate(['/threads', { board: key, url: url }])
        // this.board.emit(this.boards[0].public_key);
    }
}
