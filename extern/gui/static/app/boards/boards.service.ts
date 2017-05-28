import { Injectable } from '@angular/core';
import { Http, Response} from '@angular/http';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/do';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import 'rxjs/add/observable/throw';

import { Boards } from './boards'

@Injectable()
export class BoardsService{
    private _boardsUrl = 'http://127.0.0.1:6420/api/get_boards';

    constructor(private _http: Http){ }

    getBoards(): Observable<Boards[]>{
        return this._http.get(this._boardsUrl)
        .map((response: Response) => <Boards[]> response.json())
        .do(data => console.log('All: ' + JSON.stringify(data)))
        .catch(this.handleError);
    }

    private handleError(error: Response){
        console.error(error);
        return Observable.throw(error.json().error || 'Server error');
    }
}