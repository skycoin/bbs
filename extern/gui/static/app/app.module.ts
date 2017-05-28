import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpModule } from '@angular/http';

import { BoardsService } from './boards/boards.service';
import { AppComponent }  from './app.component';
import { BoardsListComponent } from './boards/boards-list.component';


@NgModule({
  imports: [ BrowserModule, HttpModule ],
  declarations: [ AppComponent, BoardsListComponent ],
  bootstrap: [ BoardsListComponent ],
  providers: [ BoardsService ]
})
export class AppModule { }
