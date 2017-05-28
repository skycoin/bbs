import { Component, OnInit }  from '@angular/core';

import { Boards } from './boards'
import { BoardsService } from './boards.service'

@Component({
    selector: 'boards-list',
    templateUrl: 'app/boards/boards-list.component.html'
})
export class BoardsListComponent implements OnInit{
    boardsTitle: string = 'Boards List';
    errorMessage: string;

    boards: Boards[];

    constructor(private _boardsService: BoardsService){

    }

    ngOnInit(): void{
        this._boardsService.getBoards()
            .subscribe(boards => this.boards = boards,
                    error => this.errorMessage = <any>error);
    }
}
