import { Component, OnInit, ViewEncapsulation, Output, EventEmitter } from '@angular/core';
import { ApiService } from "../../providers";
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
    test:false;
    boards: Array<Board> = [];
    constructor(private api: ApiService, private router: Router) {

    }

    ngOnInit(): void {
        this.api.getBoards().then(data => {
            this.boards = <Array<Board>>data;
        });
    }
    openThreads(ev: Event, key: string) {
        ev.stopImmediatePropagation();
        ev.stopPropagation();
        this.router.navigate(['/threads', { board: key }])
        // this.board.emit(this.boards[0].public_key);
    }
}
